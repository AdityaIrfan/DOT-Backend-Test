package auth_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/auth"
	"github.com/saas-be-usergroup/internal/core/domain/mailer"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	EmailAlreadyTaken  = errors.New("email has already been taken")
	NoCredentialsFound = errors.New("no credentials found")
	OTPNotFound        = errors.New("otp not found")
	InvalidOTP         = errors.New("invalid otp")
)

var (
	RegisterBeforeWithEmailSuccess = func(email string) string {
		return fmt.Sprintf("Thank you for registering. An email has been sent to %s. Please check your email and use the token to finish registration.", email)
	}
)

type authService struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

func NewAuthService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) ports.AuthService {
	return &authService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (a authService) RegisterBeforeWithEmail(ctx context.Context, in auth.RegisterBeforeWithEmail) (*auth.RegisterBeforeResponse, error) {
	// check email if exist
	user, err := in.ToUser().GetOneByEmail(a.db)
	if err != nil {
		a.logger.Error("failed to get user by email : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check if user is not empty and status is verified
	if !user.IsEmpty() && user.IsVerified() {
		return nil, responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(EmailAlreadyTaken.Error()))
	}

	// create or update user
	user, err = in.ToUser().Create(a.db)
	if err != nil {
		a.logger.Error("failed to create user : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// create otp
	otp, err := auth.NewOTP(user.GetUUID(), auth.GenerateNewOTP()).Create(ctx, a.redis, 180)
	if err != nil {
		a.logger.Error("failed to set otp on redis : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// send otp to email
	go func() {
		if err := mailer.Send(ctx, *mailer.NewRegisterOTPMailer(user.GetEmail(), otp.ToRegisterMail())); err != nil {
			a.logger.Error("failed to send email : ", zap.Error(err))
		}
	}()

	return otp.ToRegisterResponse(), nil
}

func (a authService) ConfirmationRegister(ctx context.Context, in auth.ConfirmationRegister) (*auth.ConfirmationResponse, error) {
	// check user by uuid
	user, err := in.ToUser().GetOneUUID(a.db)
	if err != nil {
		a.logger.Error("failed to get user by uuid : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check empty user
	if user.IsEmpty() {
		return nil, responseErr.New(fiber.StatusNotFound, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	// check otp
	otp, err := auth.GetOTPByUUID(ctx, a.redis, in.GetUUID())
	if err != nil {
		a.logger.Error("failed to get otp by uuid : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check otp is not found
	if otp.IsNotFound() {
		return nil, responseErr.New(fiber.StatusNotFound, responseErr.WithMessage(OTPNotFound.Error()))
	}

	// check invalid otp
	if !otp.IsValid(in.GetOTP()) {
		return nil, responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(InvalidOTP.Error()))
	}

	// generate session token
	sessionToken, err := auth.GenerateSessionToken(ctx, a.redis, user.GetUUID(), 1800)
	if err != nil {
		a.logger.Error("failed to create session token : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// update user status
	user.SetStatusConfirmed()
	if _, err = user.Update(a.db); err != nil {
		a.logger.Error("failed to update user : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return sessionToken.ToConfirmationResponse(), nil
}
