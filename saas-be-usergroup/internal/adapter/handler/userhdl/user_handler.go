package userhdl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/saas-be-usergroup/internal/core/domain/user"
	"github.com/saas-be-usergroup/internal/core/ports"
	responseErr "github.com/saas-be-usergroup/internal/error"
	"github.com/saas-be-usergroup/internal/response"
)

type userHandler struct {
	App         *fiber.App
	userService ports.UserService
}

func NewUserHandler(app *fiber.App, userService ports.UserService) *userHandler {
	return &userHandler{
		App:         app,
		userService: userService,
	}
}

func (u userHandler) IsEmailAvailable(c *fiber.Ctx) error {
	in := user.NewIsEmailAvailableRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := u.userService.IsEmailAvailable(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}

func (u userHandler) IsUsernameAvailable(c *fiber.Ctx) error {
	in := user.NewIsUsernameAvailableRequest()
	if err := c.BodyParser(&in); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	if err := in.Validate(); err != nil {
		return responseErr.Response(c, responseErr.New(fiber.StatusBadRequest, responseErr.WithMessage(err.Error())))
	}
	res, err := u.userService.IsUsernameAvailable(c.Context(), *in)
	if err != nil {
		return responseErr.Response(c, err)
	}
	return response.Success(c, fiber.StatusOK, response.SuccessData(res))
}
