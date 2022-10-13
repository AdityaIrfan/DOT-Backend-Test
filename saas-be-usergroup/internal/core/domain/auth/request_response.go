package auth

import (
	"github.com/google/uuid"
	"github.com/saas-be-usergroup/internal/core/domain/user"
	"strings"
)

type RegisterBeforeWithEmail struct {
	Email string `json:"email"`
}

func NewRegisterBeforeWithEmail() *RegisterBeforeWithEmail {
	return &RegisterBeforeWithEmail{}
}

func (r RegisterBeforeWithEmail) GetEmail() string {
	return r.Email
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

func NewConfirmationRegister() *ConfirmationRegister {
	return &ConfirmationRegister{}
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

type DoRegisterRequest struct {
	SessionToken    string `json:"sessionToken"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func NewDoRegisterRequest() *DoRegisterRequest {
	return &DoRegisterRequest{}
}

func (d DoRegisterRequest) GetUsername() string {
	return d.UserName
}

func (d DoRegisterRequest) GetSessionToken() string {
	return d.SessionToken
}

func (d DoRegisterRequest) GetPassword() string {
	return d.Password
}

func (d DoRegisterRequest) GetFirstName() string {
	return d.FirstName
}

func (d DoRegisterRequest) GetLastName() string {
	return d.LastName
}

func (d DoRegisterRequest) ToUpdateUser(u *user.User, hashedPassword string) *user.User {
	u.SetFirstName(d.GetFirstName())
	u.SetLastName(d.GetLastName())
	u.SetUsername(d.GetUsername())
	u.SetPassword(hashedPassword)
	return u
}

func IsPasswordMatched(password string, confirmPassword string) bool {
	return strings.EqualFold(password, confirmPassword)
}

type DoRegisterResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type DoLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewDoLoginRequest() *DoLoginRequest {
	return &DoLoginRequest{}
}

func (d DoLoginRequest) GetEmail() string {
	return d.Email
}

func (d DoLoginRequest) GetPassword() string {
	return d.Password
}

type DoLoginResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type DoRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewDoRefreshTokenRequest() *DoRefreshTokenRequest {
	return &DoRefreshTokenRequest{}
}

func (d DoRefreshTokenRequest) GetRefreshToken() string {
	return d.RefreshToken
}

type DoRefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type DoLogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewDoLogoutRequest() *DoLogoutRequest {
	return &DoLogoutRequest{}
}

func (d DoLogoutRequest) GetRefreshToken() string {
	return d.RefreshToken
}
