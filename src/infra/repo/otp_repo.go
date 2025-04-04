package repo

import (
	"TTCS/src/core/domain"
	"context"
)

type OtpRepo struct {
	*BaseRepo
}

func NewOtpRepo(baseRepo *BaseRepo) domain.OtpRepo {
	return &OtpRepo{
		BaseRepo: baseRepo,
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

func (r *OtpRepo) DeleteByEmail(ctx context.Context, email string) error {
	if err := r.db.Where("email =?", email).Delete(&domain.Otp{}).Error; err != nil {
		return r.returnError(ctx, err)
	}
	return nil
}
