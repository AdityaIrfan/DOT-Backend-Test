package auth

import "encoding/base64"

type Session struct {
	UUID string
}

func NewSession(uuid string) *Session {
	return &Session{UUID: uuid}
}

func (s *Session) IsEmpty() bool {
	return s == nil
}

func (s Session) GetUUID() string {
	return s.UUID
}

func (s Session) IsUUIDEmpty() bool {
	return s.UUID == ""
}

type SessionToken struct {
	SessionToken string
}

func NewSessionToken(sessionToken string) *SessionToken {
	return &SessionToken{SessionToken: sessionToken}
}

func (s SessionToken) ToConfirmationResponse() *ConfirmationResponse {
	return &ConfirmationResponse{SessionToken: s.SessionToken}
}

func GenerateNewSession(uuid string) string {
	return base64.StdEncoding.EncodeToString([]byte(uuid))
}
