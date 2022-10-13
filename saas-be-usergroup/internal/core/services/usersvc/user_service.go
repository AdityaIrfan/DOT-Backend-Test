package usersvc

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/user"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserService(db *gorm.DB, logger *zap.Logger) ports.UserService {
	return &userService{db: db, logger: logger}
}

func (u userService) IsEmailAvailable(ctx context.Context, in user.IsEmailAvailableRequest) (*user.AvailableResponse, error) {
	isAvailable, err := user.NewUser().IsEmailAvailable(u.db, in.GetEmail())
	if err != nil {
		u.logger.Error("failed to check email available : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return user.NewAvailableResponse(isAvailable), nil
}

func (u userService) IsUsernameAvailable(ctx context.Context, in user.IsUsernameAvailableRequest) (*user.AvailableResponse, error) {
	isAvailable, err := user.NewUser().IsUsernameAvailable(u.db, in.GetUsername())
	if err != nil {
		u.logger.Error("failed to check username available : ", zap.Error(err))
		return nil, responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return user.NewAvailableResponse(isAvailable), nil
}

func (u userService) ChangeEmailBefore(ctx context.Context, in user.ChangeEmailBeforeRequest, userAccountUUID string) error {
	return nil
}

func (u userService) ChangeEmailConfirmation(ctx context.Context, in user.ChangeEmailConfirmationRequest, userAccountUUID string) error {
	return nil
}

func (u userService) ChangePasswordRequest(ctx context.Context, userAccountUUID string) error {
	return nil
}

func (u userService) ChangePasswordConfirmation(ctx context.Context, in user.ChangePasswordConfirmationRequest, userAccountUUID string) error {
	return nil
}

func (u userService) DoChangePassword(ctx context.Context, in user.DoChangePasswordRequest, userAccountUUID string) error {
	return nil
}

func (u userService) UpdateUser(ctx context.Context, in user.UpdateRequest, userAccountUUID string) (*user.UpdateResponse, error) {
	return nil, nil
}
