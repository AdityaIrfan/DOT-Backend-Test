package auth

import "encoding/base64"

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
