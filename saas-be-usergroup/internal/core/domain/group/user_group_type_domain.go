package group

type GroupType string

const (
	GroupFree       GroupType = "free"
	GroupProject    GroupType = "project"
	GroupCommercial GroupType = "commercial"
	GroupCompany    GroupType = "company"
)

type UserGroupType struct {
	ID        uint64
	TypeName  GroupType
	MemberMin int
	MemberMax int
}

func NewUserGroupType() *UserGroupType {
	return &UserGroupType{}
}

func (u *UserGroupType) SetID(id uint64) *UserGroupType {
	u.ID = id
	return u
}

func (u *UserGroupType) IsEmpty() bool {
	return u == nil
}

func (u *UserGroupType) IsAvailableToAddOneMore(totalMember int) bool {
	return u.MemberMax >= totalMember+1
}
