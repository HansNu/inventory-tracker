package repository

import (
	"context"

	"inventory-tracker/apperr"
	"inventory-tracker/models"
)

type CategoryRepository interface {
	List(ctx context.Context) ([]models.AssetCategory, error)
	Create(ctx context.Context, category models.AssetCategory) error
	DeleteByID(ctx context.Context, id int) error
}

type categoryRepo struct {
	DB DBTX
}

func NewCategoryRepo(db DBTX) CategoryRepository {
	return &categoryRepo{DB: db}
}

func (r *categoryRepo) List(ctx context.Context) ([]models.AssetCategory, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, category_name, category_group FROM category ORDER BY category_name`)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	categories := []models.AssetCategory{}
	for rows.Next() {
		var a models.AssetCategory
		if err := rows.Scan(&a.Id, &a.CategoryName, &a.CategoryGroup); err != nil {
			return nil, translate(err)
		}
		categories = append(categories, a)
	}
	// rows.Err() reports a failure that happened mid-iteration — without it a
	// truncated result set looks identical to a complete one.
	if err := rows.Err(); err != nil {
		return nil, translate(err)
	}
	return categories, nil
}

func (r *categoryRepo) Create(ctx context.Context, category models.AssetCategory) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO category (category_name, category_group) VALUES ($1, $2)`,
		category.CategoryName, category.CategoryGroup)
	return translate(err)
}

func (r *categoryRepo) DeleteByID(ctx context.Context, id int) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM category WHERE id = $1`, id)
	if err != nil {
		return translate(err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
