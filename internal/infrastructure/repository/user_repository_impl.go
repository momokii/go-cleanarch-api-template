package repository

import (
	"context"
	"database/sql"
	"gofiber-cleanarch-test/internal/domain/entity"
	"gofiber-cleanarch-test/internal/domain/repository"
)

type UserRepositoryImpl struct{}

func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, user *entity.User) (entity.User, error) {
	var id int
	sql := "insert into users (username, password, role) values ($1, $2, $3) returning id"
	result := tx.QueryRowContext(ctx, sql, user.Username, user.Password, user.Role)

	if err := result.Scan(&id); err != nil {
		return *user, err
	}

	user.Id = int(id)

	return *user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	sql := "update users set username = $1, role = $2, updated_at = NOW() where id = $3"
	if _, err := tx.ExecContext(ctx, sql, user.Username, user.Role, user.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	sql := "delete from users where id = $1"
	if _, err := tx.ExecContext(ctx, sql, user.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) ChangePassword(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	sql := "update users set password = $1, updated_at = NOW() where id = $2"
	if _, err := tx.ExecContext(ctx, sql, user.Password, user.Id); err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, id int) (entity.User, error) {
	var user entity.User

	sql := "select id, username, password, role, created_at, updated_at, is_deleted from users where id = $1 and is_deleted = false"

	if err := tx.QueryRowContext(ctx, sql, id).Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsDeleted); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, tx *sql.Tx, username string) (entity.User, error) {
	var user entity.User

	sql := "select id, username, password, role, created_at, updated_at, is_deleted from users where username = $1 and is_deleted = false"

	if err := tx.QueryRowContext(ctx, sql, username).Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsDeleted); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindAllWithPagination(ctx context.Context, tx *sql.Tx, limit int, offset int) ([]entity.User, error) {
	var users []entity.User

	sql := "select id, username, role, created_at, updated_at from users where is_deleted=false limit $1 offset $2"
	rows, err := tx.QueryContext(ctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepositoryImpl) FindTotal(ctx context.Context, tx *sql.Tx) (int, error) {
	var total int

	sql := "select count(id) from users where is_deleted = false"
	if err := tx.QueryRowContext(ctx, sql).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}
