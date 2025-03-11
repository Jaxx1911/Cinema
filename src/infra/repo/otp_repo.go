package repo

import (
	"TTCS/src/core/domain"
	"context"
	"gorm.io/gorm"
)

type OtpRepo struct {
	*BaseRepo
	db *gorm.DB
}

func NewOtpRepo(baseRepo *BaseRepo, db *gorm.DB) domain.OtpRepo {
	return &OtpRepo{
		BaseRepo: baseRepo,
		db:       db,
	}
}

func (r *OtpRepo) Create(ctx context.Context, otp *domain.Otp) error {
	if err := r.db.Create(otp).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}

func (r *OtpRepo) GetByEmail(ctx context.Context, email string) (*domain.Otp, error) {
	var otp domain.Otp
	if err := r.db.Where("email = ?", email).First(&otp).Error; err != nil {
		return nil, r.returnError(ctx, err)
	}
	return &otp, nil
}
