package helper

import (
	"gofiber-cleanarch-test/internal/domain/entity"
	"gofiber-cleanarch-test/internal/interfaces/http/dto"
)

// for user domain response
func ToUserResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponses(users []entity.User) []dto.UserResponse {
	var usersRes []dto.UserResponse

	if users == nil {
		return []dto.UserResponse{}
	}

	for _, user := range users {
		usersRes = append(usersRes, ToUserResponse(user))
	}

	return usersRes
}
