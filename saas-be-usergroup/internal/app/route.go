package app

import (
	"github.com/go-redis/redis/v8"
	"github.com/saas-be-usergroup/internal/adapter/handler/authhdl"
	"github.com/saas-be-usergroup/internal/adapter/handler/userhdl"
	"github.com/saas-be-usergroup/internal/core/middleware"
	"github.com/saas-be-usergroup/internal/core/services/authsvc"
	"github.com/saas-be-usergroup/internal/core/services/usersvc"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handlers struct {
	Postgres *gorm.DB
	Redis    *redis.Client
	R        *fiber.App
	Logger   *zap.Logger
}

const apiVerion string = "api/v1"

func (h *Handlers) SetupRouter() {
	h.R.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	//initialize bussiness
	authService := authsvc.NewAuthService(h.Postgres, h.Redis, h.Logger)
	userService := usersvc.NewUserService(h.Postgres, h.Logger)

	//handlers initialize
	authHandler := authhdl.NewAuthHandler(h.R, authService)
	userHandler := userhdl.NewUserHandler(h.R, userService)

	// Auth
	authApi := authHandler.App.Group(apiVerion + "/auth")
	// Register
	authApi.Post("/register/before", authHandler.RegisterBeforeWithEmail)
	authApi.Post("/register/confirmation", authHandler.ConfirmationRegister)
	authApi.Post("/register/do", authHandler.DoRegister)
	// Login
	authApi.Post("/login/do", authHandler.DoLogin)
	authApi.Post("/refresh", authHandler.DoRefreshToken)
	authApi.Post("/logout", middleware.Protected(), authHandler.DoLogout)

	// User
	userApiPublic := userHandler.App.Group(apiVerion)
	userApiPublic.Post("/email/available", userHandler.IsEmailAvailable)
	userApiPublic.Post("/username/available", userHandler.IsUsernameAvailable)

}
