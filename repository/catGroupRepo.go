package repository

import (
	"context"

	"inventory-tracker/apperr"
	"inventory-tracker/models"
)

type CategoryGroupRepository interface {
	List(ctx context.Context) ([]models.CategoryGroup, error)
	Create(ctx context.Context, group models.CategoryGroup) error
	DeleteByID(ctx context.Context, id int) error
}

type catGroupRepo struct {
	DB DBTX
}

func NewCategoryGroupRepo(db DBTX) CategoryGroupRepository {
	return &catGroupRepo{DB: db}
}

func (r *catGroupRepo) List(ctx context.Context) ([]models.CategoryGroup, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, group_name FROM category_group ORDER BY group_name`)
	if err != nil {
		return nil, translate(err)
	}
	defer rows.Close()

	groups := []models.CategoryGroup{}
	for rows.Next() {
		var g models.CategoryGroup
		if err := rows.Scan(&g.Id, &g.GroupName); err != nil {
			return nil, translate(err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, translate(err)
	}
	return groups, nil
}

func (r *catGroupRepo) Create(ctx context.Context, group models.CategoryGroup) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO category_group (group_name) VALUES ($1)`, group.GroupName)
	return translate(err)
}

func (r *catGroupRepo) DeleteByID(ctx context.Context, id int) error {
	tag, err := r.DB.Exec(ctx, `DELETE FROM category_group WHERE id = $1`, id)
	if err != nil {
		return translate(err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
