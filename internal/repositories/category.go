package repositories

import (
	"context"
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/model"
	"github.com/igorschechtel/finance-tracker-backend/db/model/app_db/public/table"
)

type CategoryRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Category, error)
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]model.Category, error) {
	query := table.Category.SELECT(
		table.Category.AllColumns,
	).FROM(
		table.Category,
	).WHERE(
		table.Category.UserID.EQ(postgres.UUID(userID)),
	).ORDER_BY(
		table.Category.CreatedAt.DESC(),
	).LIMIT(int64(limit)).OFFSET(int64(offset))

	var dest []model.Category
	err := query.QueryContext(ctx, r.db, &dest)

	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	query := table.Category.INSERT(
		table.Category.UserID,
		table.Category.Name,
		table.Category.Description,
		table.Category.ColorHex,
	).VALUES(
		category.UserID,
		category.Name,
		category.Description,
		category.ColorHex,
	).RETURNING(table.Category.AllColumns)

	err := query.QueryContext(ctx, r.db, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}
