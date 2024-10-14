package helper

import "github.com/gofiber/fiber/v2"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error() string method is used to implement error interfaces
func (e AppError) Error() string {
	return e.Message
}

// --------------------- auth login error
func NewErrorAuthLoginUnauthorized() AppError {
	return AppError{
		Code:    fiber.StatusUnauthorized,
		Message: "Username or password incorrect",
	}
}

// ---------------------  user error
func NewErrorUserNotFound() AppError {
	return AppError{
		Code:    fiber.StatusNotFound,
		Message: "User not found",
	}
}

func NewErrorUserUsernameExist() AppError {
	return AppError{
		Code:    fiber.StatusBadRequest,
		Message: "Username already exist",
	}
}

func NewErrorUserPasswordIncorrect() AppError {
	return AppError{
		Code:    fiber.StatusBadRequest,
		Message: "Password incorrect",
	}
}
