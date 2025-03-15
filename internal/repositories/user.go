package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/table"
)

type UserRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	stmt := postgres.SELECT(
		table.User.AllColumns,
	).FROM(
		table.User,
	).ORDER_BY(
		table.User.ID.ASC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []model.User
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	stmt := table.User.INSERT(
		table.User.ID,
	).MODEL(
		user,
	).RETURNING(
		table.User.AllColumns,
	)

	var createdUser []model.User
	err := stmt.Query(r.db, &createdUser)
	if err != nil {
		return nil, err
	}

	return &createdUser[0], nil
}
