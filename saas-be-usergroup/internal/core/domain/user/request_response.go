package user

import (
	"strings"
)

func IsPasswordMatched(password string, confirmPassword string) bool {
	return strings.EqualFold(password, confirmPassword)
}

type UpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	UserName  *string `json:"user_name"`
}

type UpdateResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
}

type IsEmailAvailableRequest struct {
	Email string `json:"email"`
}

func NewIsEmailAvailableRequest() *IsEmailAvailableRequest {
	return &IsEmailAvailableRequest{}
}

func (i IsEmailAvailableRequest) GetEmail() string {
	return i.Email
}

type IsUsernameAvailableRequest struct {
	Username string `json:"username"`
}

func NewIsUsernameAvailableRequest() *IsUsernameAvailableRequest {
	return &IsUsernameAvailableRequest{}
}

func (i IsUsernameAvailableRequest) GetUsername() string {
	return i.Username
}

type AvailableResponse struct {
	IsAvailable bool `json:"is_available"`
}

func NewAvailableResponse(isAvailable bool) *AvailableResponse {
	return &AvailableResponse{IsAvailable: isAvailable}
}

type ChangeEmailBeforeRequest struct {
	Email string `json:"email"`
}

type ChangeEmailBeforeResponse struct {
	UUID string `json:"uuid"`
}

type ChangeEmailConfirmationRequest struct {
	OTP  string `json:"otp"`
	UUID string `json:"uuid"`
}

type ChangePasswordConfirmationRequest struct {
	OTP string `json:"otp"`
}

type DoChangePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
