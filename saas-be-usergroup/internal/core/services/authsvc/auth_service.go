package authsvc

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/auth"
	"github.com/saas-be-usergroup/internal/core/domain/mailer"
	"github.com/saas-be-usergroup/internal/core/domain/user"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	EmailAlreadyTaken    = errors.New("email has been already taken")
	NoCredentialsFound   = errors.New("no credentials found")
	OTPNotFound          = errors.New("otp not found")
	InvalidOTP           = errors.New("invalid otp")
	UsernameNotAvailable = errors.New("username not available")
	InvalidSession       = errors.New("invalid session")
	InvalidPassword      = errors.New("invalid password")
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
	user, err := user.NewUser().GetOneByEmail(a.db, in.GetEmail())
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
	user, err := user.NewUser().GetOneByUUID(a.db, in.UUID)
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

func (a authService) DoRegister(ctx context.Context, in auth.DoRegisterRequest) (*auth.DoRegisterResponse, error) {
	// check username
	isAvailable, err := user.NewUser().IsUsernameAvailable(a.db, in.GetUsername())
	if err != nil {
		a.logger.Error("failed to check username available : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if !isAvailable {
		return nil, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(UsernameNotAvailable.Error()))
	}

	// check session token
	session, err := auth.ValidateSession(ctx, a.redis, in.GetSessionToken())
	if err != nil {
		a.logger.Error("failed to validate session token : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if session.IsEmpty() || session.IsUUIDEmpty() {
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(InvalidSession.Error()))
	}

	// get user account by uuid
	userAccount, err := user.NewUser().GetOneByUUID(a.db, session.GetUUID())
	if err != nil {
		a.logger.Error("failed to get user account by uuid : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if userAccount.IsEmpty() || !userAccount.IsConfirmed() {
		return nil, responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	// generate hash password
	hashedPassword, err := auth.GeneratePassword(in.GetPassword())
	if err != nil {
		a.logger.Error("failed to generate password : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// generate JWT
	JWT, err := auth.NewJWT().Generate(ctx, a.redis, userAccount.GetUUID())
	if err != nil {
		a.logger.Error("failed to generate jwt : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// update user
	if _, err = in.ToUpdateUser(userAccount, hashedPassword).Update(a.db); err != nil {
		a.logger.Error("failed to update user : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}
	return JWT.ToDoRegisterResponse(), nil
}

func (a authService) DoLogin(ctx context.Context, in auth.DoLoginRequest) (*auth.DoLoginResponse, error) {
	// get user by email
	userAccount, err := user.NewUser().GetOneByEmail(a.db, in.GetEmail())
	if err != nil {
		a.logger.Error("failed to get user account by email : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if userAccount.IsEmpty() || !userAccount.IsVerified() {
		return nil, responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	// compare password
	if err = auth.ComparePassword(userAccount.GetPassword(), in.GetPassword()); err != nil {
		return nil, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(InvalidSession.Error()))
	}

	// generate jwt
	JWT, err := auth.NewJWT().Generate(ctx, a.redis, userAccount.GetUUID())
	if err != nil {
		a.logger.Error("failed to generate jwt : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return JWT.ToDoLoginResponse(), nil
}

func (a authService) DoRefreshToken(ctx context.Context, in auth.DoRefreshTokenRequest) (*auth.DoRefreshTokenResponse, error) {
	dataClaims, err := auth.ValidateRefreshToken(ctx, a.redis, in.GetRefreshToken())
	if err != nil {
		a.logger.Error("failed to validate refresh token : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// get user account
	userAccount, err := user.NewUser().GetOneByUUID(a.db, dataClaims.GetUUID())
	if err != nil {
		a.logger.Error("failed to get user account by email : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if userAccount.IsEmpty() || !userAccount.IsVerified() {
		return nil, responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	JWT, err := auth.NewJWT().Generate(ctx, a.redis, userAccount.GetUUID())
	if err != nil {
		a.logger.Error("failed to generate jwt : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return JWT.ToDoRefreshTokenResponse(), nil
}
func (a authService) DoLogout(ctx context.Context, in auth.DoLogoutRequest) error {
	dataClaims, err := auth.ValidateRefreshToken(ctx, a.redis, in.GetRefreshToken())
	if err != nil {
		a.logger.Error("failed to validate refresh token : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// get user account
	userAccount, err := user.NewUser().GetOneByUUID(a.db, dataClaims.GetUUID())
	if err != nil {
		a.logger.Error("failed to get user account by email : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if userAccount.IsEmpty() || !userAccount.IsVerified() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	// delete refresh token on redis
	if err = auth.DeleteRefreshToken(ctx, a.redis, dataClaims.GetJTI()); err != nil {
		a.logger.Error("failed to delete refresh token on redis : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return nil
}
