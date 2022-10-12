package user

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c CreateRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required, validation.Min(1)),
		validation.Field(&c.LastName, validation.Required, validation.Min(1)),
		validation.Field(&c.UserName, validation.Required, validation.Min(1)),
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

func (c UpdateRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Min(1)),
		validation.Field(&c.LastName, validation.Min(1)),
		validation.Field(&c.UserName, validation.Min(1)),
	)
}

func (c ChangeEmailRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
	)
}

func (c ChangePasswordRequest) Validate() error {
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
