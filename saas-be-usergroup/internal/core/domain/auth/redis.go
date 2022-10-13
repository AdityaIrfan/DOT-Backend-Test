package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func (o OTP) Create(ctx context.Context, client *redis.Client, expiredInSecond uint64) (*OTP, error) {
	expired, err := SetExpiredInSecond(expiredInSecond)
	if err != nil {
		return nil, err
	}

	o.SetDuration(expiredInSecond)

	if err = client.Set(ctx, "otp-"+o.GetUUID(), o.GetOTP(), *expired).Err(); err != nil {
		return nil, err
	}

	return &o, nil
}

func GetOTPByUUID(ctx context.Context, client *redis.Client, uuid string) (*OTP, error) {
	otp, err := client.Get(ctx, "otp-"+uuid).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return NewOTP(uuid, ""), nil
		}
		return nil, err
	}

	return NewOTP(uuid, otp), nil
}

func GenerateSessionToken(ctx context.Context, client *redis.Client, uuid string, expiredInSecond uint64) (*SessionToken, error) {
	sessionToken := GenerateNewSession(uuid)
	expired, err := SetExpiredInSecond(expiredInSecond)
	if err != nil {
		return nil, err
	}

	if err = client.SetEX(ctx, sessionToken, uuid, *expired).Err(); err != nil {
		return nil, err
	}

	return NewSessionToken(sessionToken), nil
}

func ValidateSession(ctx context.Context, client *redis.Client, sessionToken string) (*Session, error) {
	uuid, err := client.Get(ctx, sessionToken).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return NewSession(uuid), nil
}

func ValidateRefreshToken(ctx context.Context, client *redis.Client, refreshToken string) (*RefreshTokenDataClaims, error) {
	token, err := jwtParse(refreshToken, viper.GetString("jwt.refresh_secret"))
	if err != nil {
		return nil, err
	}

	jti := token.Claims.(jwt.MapClaims)["jti"].(string)
	if jti == "" {
		return nil, errors.New("invalid token")
	}

	uuid, err := client.Get(ctx, jti).Result()
	if err != nil {
		return nil, err
	}

	return NewRefreshTokenDataClaims(uuid, jti), nil
}

func DeleteRefreshToken(ctx context.Context, client *redis.Client, jti string) error {
	return client.Del(ctx, jti).Err()
}
