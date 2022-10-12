package user

import "strings"

type CreateRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func IsPasswordMatched(password string, confirmPassword string) bool {
	return strings.EqualFold(password, confirmPassword)
}

type UpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	UserName  *string `json:"user_name"`
}

type ChangeEmailRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
