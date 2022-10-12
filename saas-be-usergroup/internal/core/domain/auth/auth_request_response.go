package auth

import (
	"github.com/google/uuid"
	"github.com/saas-be-usergroup/internal/core/domain/user"
)

type RegisterBeforeWithEmail struct {
	Email string `json:"email"`
}

type RegisterBeforeResponse struct {
	UUID string `json:"uuid"`
}

func (r RegisterBeforeWithEmail) ToUser() *user.User {
	return &user.User{
		UUID:  uuid.New().String(),
		Email: r.Email,
	}
}

type ConfirmationRegister struct {
	UUID string `json:"uuid"`
	OTP  string `json:"otp"`
}

func (c ConfirmationRegister) GetUUID() string {
	return c.UUID
}

func (c ConfirmationRegister) GetOTP() string {
	return c.OTP
}

func (c ConfirmationRegister) ToUser() *user.User {
	return &user.User{UUID: c.UUID}
}

type ConfirmationResponse struct {
	SessionToken string `json:"session_token"`
}

type LoginWithEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
