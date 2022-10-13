package group

type UserGroupCreateRequest struct {
	UserGroupTypeID     uint64 `json:"user_group_type_id"`
	CustomerInvoiceData string `json:"customer_invoice_data"`
}

func (u UserGroupCreateRequest) ToUserGroup() *UserGroup {
	return &UserGroup{UserGroupTypeID: u.UserGroupTypeID, CustomerInvoiceData: u.CustomerInvoiceData}
}

type AddGroupMemberRequest struct {
	UserGroupID uint64 `json:"user_group_id"`
	Username    string `json:"username"`
}

func (a AddGroupMemberRequest) ToInGroup(userAccountID uint64) *InGroup {
	return &InGroup{
		UserGroupID:   a.UserGroupID,
		UserAccountID: userAccountID,
	}
}

type RemoveGroupMemberRequest struct {
	UserGroupID   uint64 `json:"user_group_id"`
	UserAccountID string `json:"username"`
}
