package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/clearflow-backend/db/model/app_db/public/table"
	u "github.com/igorschechtel/clearflow-backend/internal/utils"
)

type UserRepository interface {
	List(ctx context.Context, limit, offset int) ([]model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Upsert(ctx context.Context, user *model.User) (*model.User, error)
	DeleteByClerkID(ctx context.Context, clerkID string) error
	GetInternalIDByClerkID(ctx context.Context, clerkID string) (uuid.UUID, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	stmt := table.User.SELECT(
		table.User.AllColumns,
	).FROM(
		table.User,
	).ORDER_BY(
		table.User.ID.ASC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []model.User
	err := stmt.QueryContext(ctx, r.db, &dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	stmt := table.User.INSERT(
		table.User.ClerkID,
		table.User.Email,
		table.User.FirstName,
		table.User.LastName,
		table.User.ImageURL,
	).MODEL(
		user,
	).RETURNING(
		table.User.AllColumns,
	)

	var createdUser model.User
	err := stmt.QueryContext(ctx, r.db, &createdUser)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *userRepository) Upsert(ctx context.Context, user *model.User) (*model.User, error) {
	stmt := table.User.INSERT(
		table.User.ClerkID,
		table.User.Email,
		table.User.FirstName,
		table.User.LastName,
		table.User.ImageURL,
	).MODEL(
		user,
	).ON_CONFLICT(
		table.User.ClerkID,
	).DO_UPDATE(
		postgres.SET(
			table.User.Email.SET(table.User.EXCLUDED.Email),
			table.User.FirstName.SET(table.User.EXCLUDED.FirstName),
			table.User.LastName.SET(table.User.EXCLUDED.LastName),
			table.User.ImageURL.SET(table.User.EXCLUDED.ImageURL),
		),
	).RETURNING(
		table.User.AllColumns,
	)

	var upsertedUser model.User
	err := stmt.QueryContext(ctx, r.db, &upsertedUser)
	if err != nil {
		return nil, err
	}

	return &upsertedUser, nil
}

func (r *userRepository) DeleteByClerkID(ctx context.Context, clerkID string) error {
	stmt := table.User.DELETE().WHERE(table.User.ClerkID.EQ(postgres.String(clerkID)))

	_, err := stmt.ExecContext(ctx, r.db)
	return err
}

func (r *userRepository) GetInternalIDByClerkID(ctx context.Context, clerkID string) (uuid.UUID, error) {
	stmt := table.User.SELECT(table.User.ID).WHERE(table.User.ClerkID.EQ(postgres.String(clerkID)))

	var user model.User
	err := stmt.QueryContext(ctx, r.db, &user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, u.ErrNotFound
		}
		return uuid.Nil, err
	}

	return user.ID, nil
}
