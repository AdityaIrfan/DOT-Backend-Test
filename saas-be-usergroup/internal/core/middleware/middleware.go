package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
	fiberJwt "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"github.com/spf13/viper"
)

func Protected() func(*fiber.Ctx) error {
	return fiberJwt.New(fiberJwt.Config{
		SigningKey:   []byte(viper.GetString("jwt.access_secret")),
		ErrorHandler: jwtError})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	} else {
		c.Status(fiber.StatusUnauthorized)
		return responseErr.Response(c, responseErr.New(fiber.StatusUnauthorized, responseErr.WithMessage(err.Error())))
	}
}

type JWTData struct {
	UUID string
}

func NewJWTData(uuid string) *JWTData {
	return &JWTData{UUID: uuid}
}

func (j JWTData) GetUUID() string {
	return j.UUID
}

func ExportData(c context.Context) *JWTData {
	token := c.Value("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return NewJWTData(claims["uuid"].(string))
}
