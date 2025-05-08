package repo

import (
	"TTCS/src/core/domain"
	"context"

	"github.com/google/uuid"
)

type ComboRepository struct {
	*BaseRepo
}

func (c ComboRepository) Create(ctx context.Context, combo *domain.Combo) error {
	if err := c.db.WithContext(ctx).Create(combo).Error; err != nil {
		return c.returnError(ctx, err)
	}
	return nil
}

func (c ComboRepository) FindAll(ctx context.Context) ([]*domain.Combo, error) {
	var combos []*domain.Combo
	if err := c.db.WithContext(ctx).Order("price ASC").Find(&combos).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return combos, nil
}

func (c ComboRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Combo, error) {
	var combo domain.Combo
	if err := c.db.WithContext(ctx).First(&combo, id).Error; err != nil {
		return nil, c.returnError(ctx, err)
	}
	return &combo, nil
}

func (c ComboRepository) Update(ctx context.Context, combo *domain.Combo) error {
	if err := c.db.WithContext(ctx).Save(combo).Error; err != nil {
		return c.returnError(ctx, err)
	}
	return nil
}

// Delete performs a soft delete on the combo by setting the DeletedAt timestamp
func (c ComboRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Using Delete() without Unscoped() will perform a soft delete
	if err := c.db.WithContext(ctx).Delete(&domain.Combo{}, id).Error; err != nil {
		return c.returnError(ctx, err)
	}
	return nil
}

func NewComboRepo(baseRepo *BaseRepo) domain.ComboRepository {
	return ComboRepository{BaseRepo: baseRepo}
}
