package group

import (
	"time"
)

type UserGroup struct {
	ID                  uint64
	UserGroupTypeID     uint64
	CustomerInvoiceData string
	InsertTs            time.Time
}

func NewUserGroup() *UserGroup {
	return &UserGroup{}
}

func (u *UserGroup) SetInsertTS() {
	u.InsertTs = time.Now()
}

func (u UserGroup) GetUserGroupTypeID() uint64 {
	return u.UserGroupTypeID
}

func (u UserGroup) ToInGroup(UserAccountID uint64) *InGroup {
	return &InGroup{
		UserGroupID:   u.ID,
		UserAccountID: UserAccountID,
		TimeAdded:     time.Now(),
		GroupAdmin:    true,
		Creator:       true,
	}
}

func (u *UserGroup) IsEmpty() bool {
	return u == nil
}
