package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"time"
)

var (
	claimsIssuerAccess  string            = "SaaS-JWT-Access"
	claimsIssuerRefresh string            = "SaaS-JWT-Refresh"
	jwtExpiresAt        int64             = time.Now().Add(time.Duration(1) * time.Hour * 24 * 30).Unix() // 1 month
	jwtSigningMethod    jwt.SigningMethod = jwt.SigningMethodHS256
)

type JWT struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) setAccessToken(accessToken string) {
	j.RefreshToken = accessToken
}

func (j *JWT) setRefreshToken(refreshToken string) {
	j.RefreshToken = refreshToken
}

func (j *JWT) setExpiresIn(expiresIn int64) {
	j.ExpiresIn = expiresIn
}

func (g JWT) ToDoRegisterResponse() *DoRegisterResponse {
	return &DoRegisterResponse{
		RefreshToken: g.RefreshToken,
		AccessToken:  g.AccessToken,
		ExpiresIn:    g.ExpiresIn,
	}
}

func (g JWT) ToDoLoginResponse() *DoLoginResponse {
	return &DoLoginResponse{
		RefreshToken: g.RefreshToken,
		AccessToken:  g.AccessToken,
		ExpiresIn:    g.ExpiresIn,
	}
}

func (g JWT) ToDoRefreshTokenResponse() *DoRefreshTokenResponse {
	return &DoRefreshTokenResponse{
		RefreshToken: g.RefreshToken,
		AccessToken:  g.AccessToken,
		ExpiresIn:    g.ExpiresIn,
	}
}

type JWTClaims struct {
	jwt.StandardClaims
	UUID string `json:"uuid"`
}

func newJWTClaimsAccess(uuid string) *JWTClaims {
	claims := &JWTClaims{}
	claims.setStandardClaimsAccess()
	claims.setUUID(uuid)
	return claims
}

func newJWTClaimsRefresh(uuid string) *JWTClaims {
	claims := &JWTClaims{}
	claims.setStandardClaimsAccess()
	claims.setUUID(uuid)
	return claims
}

func (j *JWTClaims) setStandardClaimsAccess() {
	j.StandardClaims = jwt.StandardClaims{
		Issuer:    claimsIssuerAccess,
		ExpiresAt: jwtExpiresAt,
		Id:        uuid.New().String()}
}

func (j *JWTClaims) setStandardClaimsRefresh() {
	j.StandardClaims = jwt.StandardClaims{
		Issuer: claimsIssuerRefresh,
		Id:     uuid.New().String()}
}

func (j *JWTClaims) setUUID(uuid string) {
	j.UUID = uuid
}

func (j *JWT) Generate(ctx context.Context, client *redis.Client, uuid string) (*JWT, error) {
	accessClaims := *newJWTClaimsAccess(uuid)
	accessToken, err := jwt.NewWithClaims(jwtSigningMethod, accessClaims).SignedString([]byte(viper.GetString("jwt.access_secret")))
	if err != nil {
		return nil, err
	}

	refreshClaims := *newJWTClaimsRefresh(uuid)
	refreshToken, err := jwt.NewWithClaims(jwtSigningMethod, refreshClaims).SignedString([]byte(viper.GetString("jwt.refresh_secret")))
	if err != nil {
		return nil, err
	}

	if err = client.Set(ctx, refreshClaims.StandardClaims.Id, refreshClaims.UUID, 0).Err(); err != nil {
		return nil, err
	}

	j.setAccessToken(accessToken)
	j.setRefreshToken(refreshToken)
	j.setExpiresIn(accessClaims.ExpiresAt)

	return j, nil
}

type RefreshTokenDataClaims struct {
	UUID string
	JTI  string
}

func NewRefreshTokenDataClaims(uuid string, jti string) *RefreshTokenDataClaims {
	return &RefreshTokenDataClaims{
		UUID: uuid,
		JTI:  jti,
	}
}

func (r RefreshTokenDataClaims) GetUUID() string {
	return r.UUID
}

func (r RefreshTokenDataClaims) GetJTI() string {
	return r.JTI
}

func jwtParse(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if method, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		} else if method != jwtSigningMethod {
			return nil, errors.New("invalid token")
		}
		return []byte(secret), nil
	})
}
