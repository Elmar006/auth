package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	//"auth/service/internal/db"
	"auth/service/internal/logger"
	"auth/service/internal/model"
)

var (
	log = logger.L()
)

type UserRepo interface {
	Create(ctx context.Context, user model.Model) (int64, error)
	GetByEmail(ctx context.Context, email string) (*model.Model, error)
	GetByID(ctx context.Context, id int64) (*model.Model, error)
}

type userDB struct {
	*sql.DB
}

func NewUserRepo(data *sql.DB) UserRepo {
	return &userDB{data}
}

func (u *userDB) Create(ctx context.Context, user model.Model) (int64, error) {
	query := `INSERT INTO users (name, email, password_hash, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id`

	user.CreatedAt = time.Now()

	var id int64
	err := u.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}

	return id, nil
}

func (u *userDB) GetByEmail(ctx context.Context, email string) (*model.Model, error) {
	task := &model.Model{}
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = $1`

	err := u.QueryRowContext(ctx, query, email).Scan(
		&task.ID, &task.Name, &task.Email, &task.Password, &task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return task, nil
}

func (u *userDB) GetByID(ctx context.Context, id int64) (*model.Model, error) {
	task := &model.Model{}
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE id = $1`

	err := u.QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Name, &task.Email, &task.Password, &task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return task, nil
}
