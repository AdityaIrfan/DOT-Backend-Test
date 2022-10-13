package group

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (u UserGroupCreateRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserGroupTypeID, validation.Required),
	)
}

func (u AddGroupMemberRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserGroupID, validation.Required),
		validation.Field(&u.Username, validation.Required),
	)
}

func (u RemoveGroupMemberRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserGroupID, validation.Required),
		validation.Field(&u.UserAccountID, validation.Required),
	)
}
