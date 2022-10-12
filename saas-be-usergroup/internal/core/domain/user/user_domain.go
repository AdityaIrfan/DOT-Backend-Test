package user

import "time"

type UserStatus string

const (
	UserNew       UserStatus = "new"
	UserConfirmed            = "confirmed"
	UserVerifieed            = "verified"
)

type User struct {
	ID               uint64
	UUID             string
	FirstName        string
	LastName         string
	UserName         string
	Password         string
	Email            string
	Status           UserStatus
	ConfirmationTime time.Time
	InsertTs         time.Time
}

func (u *User) IsEmpty() bool {
	return u == nil
}

func (u User) IsNew() bool {
	return u.Status == UserNew
}

func (u User) IsConfirmed() bool {
	return u.Status == UserConfirmed
}

func (u User) IsVerified() bool {
	return u.Status == UserVerifieed
}

func (u User) GetUUID() string {
	return u.UUID
}

func (u User) GetEmail() string {
	return u.Email
}

func (u *User) SetStatusConfirmed() {
	u.Status = UserConfirmed
}

func (u *User) SetStatusVerified() {
	u.Status = UserVerifieed
}

func (u *User) ToTransformer() *Transformer {
	return &Transformer{
		UUID:      u.UUID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UserName:  u.UserName,
		Email:     u.Email,
	}
}
