package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gofiber-cleanarch-test/internal/domain/entity"
	"gofiber-cleanarch-test/internal/domain/repository"
	"gofiber-cleanarch-test/internal/interfaces/http/dto"
	"gofiber-cleanarch-test/pkg/helper"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAllWithPagination(ctx context.Context, limit int, offset int) (dto.PaginationData, error)
	FindById(ctx context.Context, Id int) (dto.UserResponse, error)
	FindByUsername(ctx context.Context, username string) (dto.UserResponse, error)
	Create(ctx context.Context, req *dto.UserCreate) (dto.UserResponse, error)
	Update(ctx context.Context, req *dto.UserUpdate) error
	ChangePassword(ctx context.Context, req *dto.UserChangePassword) error
	Delete(ctx context.Context, Id int) error
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserService(userRepository repository.UserRepository, db *sql.DB) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (s *UserServiceImpl) FindAllWithPagination(ctx context.Context, limit int, offset int) (dto.PaginationData, error) {
	res, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {
		users, err := s.UserRepository.FindAllWithPagination(ctx, tx, limit, offset)
		if err != nil {
			return dto.PaginationData{}, err
		}

		totalData, err := s.UserRepository.FindTotal(ctx, tx)
		if err != nil {
			return dto.PaginationData{}, err
		}

		return dto.PaginationData{
			TotalData: totalData,
			Data:      helper.ToUserResponses(users),
		}, nil
	})

	return res.(dto.PaginationData), err
}

func (s *UserServiceImpl) FindById(ctx context.Context, Id int) (dto.UserResponse, error) {
	res, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {
		user, err := s.UserRepository.FindByID(ctx, tx, Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return dto.UserResponse{}, helper.NewErrorUserNotFound()
			}

			return dto.UserResponse{}, err
		}

		return helper.ToUserResponse(user), nil
	})

	return res.(dto.UserResponse), err
}

func (s *UserServiceImpl) FindByUsername(ctx context.Context, username string) (dto.UserResponse, error) {
	res, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {
		user, err := s.UserRepository.FindByUsername(ctx, tx, username)
		if err != nil {
			if err == sql.ErrNoRows {
				return dto.UserResponse{}, helper.NewErrorUserNotFound()
			}

			return dto.UserResponse{}, err
		}

		return helper.ToUserResponse(user), nil
	})

	return res.(dto.UserResponse), err
}

func (s *UserServiceImpl) Create(ctx context.Context, req *dto.UserCreate) (dto.UserResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return dto.UserResponse{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("rollback error: %v (original error: %w)", rbErr, err)
			}
		}
	}()

	user := entity.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	// check username if available
	user_check, err := s.UserRepository.FindByUsername(ctx, tx, user.Username)
	if (err != nil) && (err != sql.ErrNoRows) {

		return dto.UserResponse{}, err
	}

	if user_check.Id != 0 {
		return dto.UserResponse{}, helper.NewErrorUserUsernameExist()
	}

	// hashed password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 16)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user.Password = string(hashedPass)

	// save user
	user, err = s.UserRepository.Save(ctx, tx, &user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if err = tx.Commit(); err != nil {
		return dto.UserResponse{}, errors.New("error commit transaction")
	}

	return helper.ToUserResponse(user), nil
}

func (s *UserServiceImpl) Update(ctx context.Context, req *dto.UserUpdate) error {
	_, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {

		// check user by id
		user, err := s.UserRepository.FindByID(ctx, tx, req.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, helper.NewErrorUserNotFound()
			}

			return nil, err
		}

		// check username if exist or not
		username_check, err := s.UserRepository.FindByUsername(ctx, tx, req.Username)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}

		// check if avail but not the same user
		if (username_check.Id != 0) && (username_check.Id != user.Id) {
			return nil, helper.NewErrorUserUsernameExist()
		}

		user_update := entity.User{
			Id:       req.Id,
			Username: req.Username,
			Role:     req.Role,
		}

		// update user data
		if err = s.UserRepository.Update(ctx, tx, &user_update); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, req *dto.UserChangePassword) error {
	_, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {
		// check user by id
		user, err := s.UserRepository.FindByID(ctx, tx, req.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, helper.NewErrorUserNotFound()
			}

			return nil, err
		}

		// compare password
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			return nil, helper.NewErrorUserPasswordIncorrect()
		}

		// hash new pass
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), 16)
		if err != nil {
			return nil, err
		}

		user.Password = string(hashedPass)

		// update user
		if err = s.UserRepository.ChangePassword(ctx, tx, &user); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

func (s *UserServiceImpl) Delete(ctx context.Context, Id int) error {
	_, err := helper.WithTransaction(ctx, s.DB, func(tx *sql.Tx) (interface{}, error) {
		// check user by id
		user, err := s.UserRepository.FindByID(ctx, tx, Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, helper.NewErrorUserNotFound()
			}

			return nil, err
		}

		// delete user
		if err = s.UserRepository.Delete(ctx, tx, &user); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
