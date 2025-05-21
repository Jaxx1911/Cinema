package repo

import (
	"TTCS/src/core/domain"
	"TTCS/src/present/httpui/request"
	"context"

	"github.com/google/uuid"
)

type UserRepo struct {
	*BaseRepo
}

func NewUserRepo(baseRepo *BaseRepo) domain.UserRepo {
	return &UserRepo{
		BaseRepo: baseRepo,
	}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user, nil
}

func (r *UserRepo) GetList(ctx context.Context, page request.GetListUser) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	// Apply filters
	if page.Role != "" && page.Role != "all" {
		query = query.Where("role = ?", page.Role)
	}
	if page.Name != "" {
		query = query.Where("name LIKE ?", "%"+page.Name+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	limit, offset := r.toLimitOffset(ctx, page.Page)
	if err := query.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, r.returnError(ctx, err)
	}

	return users, total, nil
}

func (r *UserRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).First(user, id).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user, nil
}
func (r *UserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user, nil
}

func (r *UserRepo) GetPaymentsById(ctx context.Context, id uuid.UUID) ([]domain.Payment, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Preload("Payments").Where("id = ?", id).First(user).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user.Payments, nil
}

func (r *UserRepo) GetOrdersById(ctx context.Context, id uuid.UUID) ([]domain.Order, error) {
	user := &domain.User{}
	if err := r.db.WithContext(ctx).Preload("Orders").Where("id = ?", id).First(user).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return user.Orders, nil
}

func (r *UserRepo) Delete(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Delete(user).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}
