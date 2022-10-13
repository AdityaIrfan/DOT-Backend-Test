package user

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c UpdateRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Min(1)),
		validation.Field(&c.LastName, validation.Min(1)),
		validation.Field(&c.UserName, validation.Min(1)),
	)
}

func (i IsUsernameAvailableRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Username, validation.Required),
	)
}

func (i IsEmailAvailableRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Email, validation.Required, is.Email),
	)
}

func (c ChangeEmailBeforeRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

func (c ChangeEmailConfirmationRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.UUID, validation.Required, is.Email),
		validation.Field(&c.OTP, validation.Required, is.Email),
	)
}

func (c ChangePasswordConfirmationRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OTP, validation.Required),
	)
}

func (c DoChangePasswordRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Password, validation.Required, validation.Min(8)),
		validation.Field(&c.ConfirmPassword, validation.Required, validation.Min(8)),
	); err != nil {
		return err
	}

	if IsPasswordMatched(c.Password, c.ConfirmPassword) {
		return errors.New("Password does not match")
	}

	return nil
}
