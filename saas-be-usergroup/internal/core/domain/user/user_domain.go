package user

import (
	"fmt"
	"time"
)

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

func (u *User) SetFirstName(firstName string) {
	u.FirstName = firstName
}

func (u *User) SetLastName(lastName string) {
	u.LastName = lastName
}

func (u *User) SetUsername(userName string) {
	u.UserName = userName
}

func (u *User) SetPassword(hashedPassword string) {
	u.Password = hashedPassword
}

func NewUser() *User {
	return &User{}
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

func (u User) GetName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u User) GetID() uint64 {
	return u.ID
}

func (u User) GetPassword() string {
	return u.Password
}
