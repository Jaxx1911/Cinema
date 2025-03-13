package repo

import (
	"TTCS/src/common/fault"
	"TTCS/src/infra/cache"
	"TTCS/src/present/httpui/request"
	"context"
	"errors"
	"gorm.io/gorm"
)

type BaseRepo struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func NewBaseRepo(db *gorm.DB, cache *cache.RedisCache) *BaseRepo {
	return &BaseRepo{
		db:    db,
		cache: cache,
	}
}

func (b *BaseRepo) returnError(ctx context.Context, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fault.Wrapf(err, "[%v] record not found", "DB").SetTag(fault.TagNotFound)
	}
	return fault.Wrapf(err, "internal")
}

func (b *BaseRepo) toLimitOffset(ctx context.Context, page request.Page) (int, int) {
	offset := page.Limit * (page.Page - 1)
	return page.Limit, offset
}
