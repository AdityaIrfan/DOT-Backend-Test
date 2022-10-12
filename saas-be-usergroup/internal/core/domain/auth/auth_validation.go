package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c RegisterBeforeWithEmail) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

func (c ConfirmationRegister) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required),
		validation.Field(&c.OTP, validation.Required),
	)
}

func (c LoginWithEmail) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required),
	)
}
