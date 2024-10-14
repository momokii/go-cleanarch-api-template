package service

import (
	"context"
	"database/sql"
	"gofiber-cleanarch-test/internal/domain/repository"
	"gofiber-cleanarch-test/internal/interfaces/http/dto"
	"gofiber-cleanarch-test/pkg/helper"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginUser(ctx context.Context, req *dto.LoginInput) (dto.LoginResponse, error)
}

type AuthServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewAuthService(userRepository repository.UserRepository, db *sql.DB) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (s *AuthServiceImpl) LoginUser(ctx context.Context, req *dto.LoginInput) (dto.LoginResponse, error) {
	res, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {

		// check user username
		user, err := s.UserRepository.FindByUsername(ctx, tx, req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				return dto.LoginResponse{}, helper.NewErrorAuthLoginUnauthorized()
			}

			return dto.LoginResponse{}, err
		}

		// check password
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return dto.LoginResponse{}, helper.NewErrorAuthLoginUnauthorized()
		}

		// create token
		sign := jwt.New(jwt.SigningMethodHS256)
		claims := sign.Claims.(jwt.MapClaims)
		claims["id"] = user.Id
		claims["exp"] = time.Now().Add(time.Hour * 8).Unix()

		token, err := sign.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return dto.LoginResponse{}, err
		}

		return dto.LoginResponse{Token: token}, nil
	})

	return res.(dto.LoginResponse), err
}
