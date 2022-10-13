package groupsvc

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/group"
	"github.com/saas-be-usergroup/internal/core/domain/user"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	NoCredentialsFound          = errors.New("no credentials found")
	UserStatusNotVerified       = errors.New("finish your registration for getting verified status")
	AdminNotFound               = errors.New("admin not found")
	MemberNotFound              = errors.New("member not found")
	AdminStatusNotVerified      = errors.New("admin have not completed the registration for getting verified status")
	MemberStatusNotVerified     = errors.New("member have not completed the registration for getting verified status")
	AdminHaveNotJoinedTheGroup  = errors.New("admin have not joined the group")
	MemberHaveNotJoinedTheGroup = errors.New("member have not joined the group")
	NotAdmin                    = errors.New("you are not admin")
	UnauthorizeToRemoveCreator  = errors.New("can not remove creator of this group")
)

type groupService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewGroupService(db *gorm.DB, logger *zap.Logger) ports.GroupService {
	return &groupService{
		db:     db,
		logger: logger,
	}
}

func (g groupService) CreateGroup(ctx context.Context, in group.UserGroupCreateRequest, creatorUUID string) error {
	// get user by uuid
	user, err := user.NewUser().GetOneByUUID(g.db, creatorUUID)
	if err != nil {
		g.logger.Error("failed to get user by uuid : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check empty user
	if user.IsEmpty() {
		return responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(NoCredentialsFound.Error()))
	}

	// check user status, user status should be verified
	if !user.IsVerified() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(UserStatusNotVerified.Error()))
	}

	// get in group
	if _, err = in.ToUserGroup().Create(g.db, user.GetID()); err != nil {
		g.logger.Error("failed to create group : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return nil
}

func (g groupService) AddGroupMember(ctx context.Context, in group.AddGroupMemberRequest, adminUUID string) error {
	// get admin by uuid
	admin, err := user.NewUser().GetOneByUUID(g.db, adminUUID)
	if err != nil {
		g.logger.Error("failed to get user admin by uuid : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check empty user
	if admin.IsEmpty() {
		return responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(AdminNotFound.Error()))
	}

	// check user status, user status should be verified
	if !admin.IsVerified() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(AdminStatusNotVerified.Error()))
	}

	// get user who want to add by uuid
	member, err := user.NewUser().GetOneByUUID(g.db, adminUUID)
	if err != nil {
		g.logger.Error("failed to get user member by uuid : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check empty user
	if member.IsEmpty() {
		return responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(MemberNotFound.Error()))
	}

	// check user status, user status should be verified
	if !member.IsVerified() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(MemberStatusNotVerified.Error()))
	}

	inGroup, err := group.NewInGroup().GetOneByUserGroupIDAndUserAccountID(g.db, in.UserGroupID, admin.ID)
	if err != nil {
		g.logger.Error("failed to get in_group by group id and user id : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if inGroup.IsEmpty() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(AdminHaveNotJoinedTheGroup.Error()))
	}

	if !inGroup.IsAdmin() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NotAdmin.Error()))
	}

	if isErrorInternal, err := inGroup.IsAvailableToAddOneMore(g.db); err != nil {
		if isErrorInternal {
			g.logger.Error("failed to check IsAvailableToAddOneMore : ", zap.Error(err))
			return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
		}
		return responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(err.Error()))
	}

	if _, err = in.ToInGroup(member.GetID()).AddMember(g.db, admin.GetID()); err != nil {
		g.logger.Error("failed to add member to group : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return nil
}

func (g groupService) RemoveGroupMember(ctx context.Context, in group.RemoveGroupMemberRequest, adminUUID string) error {
	// get admin by uuid
	admin, err := user.NewUser().GetOneByUUID(g.db, adminUUID)
	if err != nil {
		g.logger.Error("failed to get user admin by uuid : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	// check empty user
	if admin.IsEmpty() {
		return responseErr.New(fiber.StatusUnprocessableEntity, responseErr.WithMessage(AdminNotFound.Error()))
	}

	// check user status, user status should be verified
	if !admin.IsVerified() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(AdminStatusNotVerified.Error()))
	}

	adminInGroup, err := group.NewInGroup().GetOneByUserGroupIDAndUserAccountID(g.db, in.UserGroupID, admin.ID)
	if err != nil {
		g.logger.Error("failed to get admin in_group by group id and user id : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if adminInGroup.IsEmpty() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(AdminHaveNotJoinedTheGroup.Error()))
	}

	if !adminInGroup.IsAdmin() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(NotAdmin.Error()))
	}

	memberInGroup, err := group.NewInGroup().GetOneByUserGroupIDAndUserAccountID(g.db, in.UserGroupID, admin.ID)
	if err != nil {
		g.logger.Error("failed to get member in_group by group id and user id : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	if memberInGroup.IsEmpty() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(MemberHaveNotJoinedTheGroup.Error()))
	}

	if memberInGroup.IsCreator() {
		return responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(UnauthorizeToRemoveCreator.Error()))
	}

	if _, err = memberInGroup.RemoveMember(g.db, admin.GetID()); err != nil {
		g.logger.Error("failed to remove member from in_group : ", zap.Error(err))
		return responseErr.New(fiber.StatusInternalServerError, responseErr.WithMessage(responseErr.ErrInternalServer.Error()))
	}

	return nil
}
