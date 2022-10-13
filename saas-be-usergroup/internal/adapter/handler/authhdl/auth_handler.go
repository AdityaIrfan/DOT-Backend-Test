package authhdl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/auth"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"github.com/saas-be-usergroup/internal/response"
)

type authHandler struct {
	App         *fiber.App
	authService ports.AuthService
}

func NewAuthHandler(app *fiber.App, authService ports.AuthService) *authHandler {
	return &authHandler{
		App:         app,
		authService: authService,
	}
}

func (a authHandler) RegisterBeforeWithEmail(c *fiber.Ctx) error {
	in := auth.NewRegisterBeforeWithEmail()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := a.authService.RegisterBeforeWithEmail(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (a authHandler) ConfirmationRegister(c *fiber.Ctx) error {
	in := auth.NewConfirmationRegister()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := a.authService.ConfirmationRegister(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (a authHandler) DoRegister(c *fiber.Ctx) error {
	in := auth.NewDoRegisterRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := a.authService.DoRegister(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (a authHandler) DoLogin(c *fiber.Ctx) error {
	in := auth.NewDoLoginRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := a.authService.DoLogin(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (a authHandler) DoRefreshToken(c *fiber.Ctx) error {
	in := auth.NewDoRefreshTokenRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := a.authService.DoRefreshToken(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (a authHandler) DoLogout(c *fiber.Ctx) error {
	in := auth.NewDoLogoutRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := a.authService.DoLogout(c.Context(), *in); err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(map[string]interface{}{
		"message": "logout success",
	}))
}
