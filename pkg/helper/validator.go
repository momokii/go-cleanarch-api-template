package helper

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})

	return validate
}

func ValidateStruct(s interface{}) error {
	return GetValidator().Struct(s)
}
