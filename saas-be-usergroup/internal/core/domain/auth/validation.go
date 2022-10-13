package auth

import (
	"errors"
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

func (c DoLoginRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required),
	)
}

func (c DoRegisterRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.SessionToken, validation.Required),
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

func (c DoRefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.RefreshToken, validation.Required),
	)
}

func (c DoLogoutRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.RefreshToken, validation.Required),
	)
}
