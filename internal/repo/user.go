package repo

import (
	"context"

	"github.com/gingray/go-template/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
}

type User struct {
	pgPool *pgxpool.Pool
}

func NewUserRepo(pgPool *pgxpool.Pool) *User {
	return &User{pgPool: pgPool}
}

func (u *User) GetUsers(ctx context.Context) ([]entity.User, error) {
	res, err := u.pgPool.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users, err := pgx.CollectRows(res, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}
