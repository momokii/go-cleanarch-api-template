package controllers

import (
	"gofiber-cleanarch-test/internal/interfaces/http/dto"
	"gofiber-cleanarch-test/internal/service"
	"gofiber-cleanarch-test/pkg/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	loginInput := new(dto.LoginInput)
	if err := c.BodyParser(loginInput); err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := helper.ValidateStruct(loginInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				return helper.RespondError(c, fiber.StatusBadRequest, "Username is required")
			case "Password":
				return helper.RespondError(c, fiber.StatusBadRequest, "Password is required")
			default:
				return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
			}
		}
	}

	token, err := h.authService.LoginUser(c.Context(), loginInput)
	if err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	return helper.RespondWithData(c, fiber.StatusOK, "success login", fiber.Map{
		"token":        token.Token,
		"token_type":   "Bearer",
		"expired_time": "8h",
	})
}
