package repository

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/table"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	stmt := postgres.SELECT(
		table.User.AllColumns,
	).FROM(
		table.User,
	).WHERE(
		table.User.ID.EQ(postgres.UUID(uid)),
	)

	var dest []model.User
	err = stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}

	if len(dest) == 0 {
		return nil, nil // Not found
	}

	return &dest[0], nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	stmt := postgres.SELECT(
		table.User.AllColumns,
	).FROM(
		table.User,
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []model.User
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	newUser := model.User{
		ID: user.ID,
	}

	table.User.INSERT(
		table.User.ID,
	).MODEL(
		newUser,
	)

	return nil
}
