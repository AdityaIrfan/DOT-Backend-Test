package ports

import (
	"context"
	"github.com/saas-be-usergroup/internal/core/domain/auth"
	"github.com/saas-be-usergroup/internal/core/domain/group"
	"github.com/saas-be-usergroup/internal/core/domain/user"
)

type (
	AuthService interface {
		RegisterBeforeWithEmail(ctx context.Context, in auth.RegisterBeforeWithEmail) (*auth.RegisterBeforeResponse, error)
		ConfirmationRegister(ctx context.Context, in auth.ConfirmationRegister) (*auth.ConfirmationResponse, error)
		DoRegister(ctx context.Context, in auth.DoRegisterRequest) (*auth.DoRegisterResponse, error)
		DoLogin(ctx context.Context, in auth.DoLoginRequest) (*auth.DoLoginResponse, error)
		DoRefreshToken(ctx context.Context, in auth.DoRefreshTokenRequest) (*auth.DoRefreshTokenResponse, error)
		DoLogout(ctx context.Context, in auth.DoLogoutRequest) error
	}

	GroupService interface {
		CreateGroup(ctx context.Context, in group.UserGroupCreateRequest, creatorUUID string) error
		AddGroupMember(ctx context.Context, in group.AddGroupMemberRequest, adminUUID string) error
		RemoveGroupMember(ctx context.Context, in group.RemoveGroupMemberRequest, adminUUID string) error
	}

	UserService interface {
		IsEmailAvailable(ctx context.Context, in user.IsEmailAvailableRequest) (*user.AvailableResponse, error)
		IsUsernameAvailable(ctx context.Context, in user.IsUsernameAvailableRequest) (*user.AvailableResponse, error)
		ChangeEmailBefore(ctx context.Context, in user.ChangeEmailBeforeRequest, userAccountUUID string) error
		ChangeEmailConfirmation(ctx context.Context, in user.ChangeEmailConfirmationRequest, userAccountUUID string) error
		ChangePasswordRequest(ctx context.Context, userAccountUUID string) error
		ChangePasswordConfirmation(ctx context.Context, in user.ChangePasswordConfirmationRequest, userAccountUUID string) error
		DoChangePassword(ctx context.Context, in user.DoChangePasswordRequest, userAccountUUID string) error
		UpdateUser(ctx context.Context, in user.UpdateRequest, userAccountUUID string) (*user.UpdateResponse, error)
	}
)
