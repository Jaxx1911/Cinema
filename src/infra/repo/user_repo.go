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

func (r *UserRepo) GetList(ctx context.Context, page request.Page) ([]*domain.User, error) {
	var users []*domain.User
	limit, offset := r.toLimitOffset(ctx, page)
	if err := r.db.WithContext(ctx).Find(&users).Limit(limit).Offset(offset).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return users, nil
}

func (r *UserRepo) GetById(ctx context.Context, id string) (*domain.User, error) {
	userId, err := uuid.Parse(id)
	if err != nil {
		return nil, r.returnError(ctx, err)
	}
	user := &domain.User{}
	if err := r.db.WithContext(ctx).First(user, userId).Error; err != nil {
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
