package controllers

import (
	"gofiber-cleanarch-test/internal/interfaces/http/dto"
	"gofiber-cleanarch-test/internal/service"
	"gofiber-cleanarch-test/pkg/helper"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (h *UserController) GetAllUsers(c *fiber.Ctx) error {
	per_page, err := strconv.Atoi(c.Query("per_page", "10"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid per_page query value")
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid page query value")
	}
	offset := (page - 1) * per_page

	users, err := h.userService.FindAllWithPagination(c.Context(), per_page, offset)
	if err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondWithPagination(c, fiber.StatusOK, "success get users data", users.TotalData, page, per_page, "users", users.Data)
}

func (h *UserController) GetUserById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid user id")
	}

	user, err := h.userService.FindById(c.Context(), id)
	if err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondWithData(c, fiber.StatusOK, "success get user data", user)
}

func (h *UserController) CreateUser(c *fiber.Ctx) error {
	userInput := new(dto.UserCreate)
	if err := c.BodyParser(userInput); err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := helper.ValidateStruct(userInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				return helper.RespondError(c, fiber.StatusBadRequest, "Username minimal 5 karakter, merupakan alphanumerik")
			case "Password":
				return helper.RespondError(c, fiber.StatusBadRequest, "Password minimal 6 karakter, mengandung angka dan huruf besar")
			default:
				return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
			}
		}
	}

	if _, err := h.userService.Create(c.Context(), userInput); err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondMessage(c, fiber.StatusOK, "success create user")
}

func (h *UserController) EditUserPassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid user id")
	}

	userInput := new(dto.UserChangePassword)
	if err = c.BodyParser(userInput); err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
	}
	userInput.Id = id

	if err = helper.ValidateStruct(userInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return helper.RespondError(c, fiber.StatusBadRequest, "Id tidak boleh kosong")
			case "OldPassword":
				return helper.RespondError(c, fiber.StatusBadRequest, "Old Password tidak boleh kosong")
			case "Password":
				return helper.RespondError(c, fiber.StatusBadRequest, "Password minimal 6 karakter, mengandung angka dan huruf besar")
			default:
				return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
			}
		}
	}

	if err = h.userService.ChangePassword(c.Context(), userInput); err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondMessage(c, fiber.StatusOK, "success edit user password")
}

func (h *UserController) EditUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid user id")
	}

	userInput := new(dto.UserUpdate)
	if err = c.BodyParser(userInput); err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
	}
	userInput.Id = id

	if err = helper.ValidateStruct(userInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return helper.RespondError(c, fiber.StatusBadRequest, "Id tidak boleh kosong")
			case "Username":
				return helper.RespondError(c, fiber.StatusBadRequest, "Username minimal 5 karakter, merupakan alphanumerik")
			default:
				return helper.RespondError(c, fiber.StatusBadRequest, err.Error())
			}
		}
	}

	if err = h.userService.Update(c.Context(), userInput); err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondMessage(c, fiber.StatusOK, "success edit user")
}

func (h *UserController) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.RespondError(c, fiber.StatusBadRequest, "Invalid user id")
	}

	if err = h.userService.Delete(c.Context(), id); err != nil {
		if e, ok := err.(helper.AppError); ok {
			return helper.RespondError(c, e.Code, e.Message)
		}
		return helper.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.RespondMessage(c, fiber.StatusOK, "success delete user")
}
