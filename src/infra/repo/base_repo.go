package repo

import (
	"TTCS/src/common"
	"TTCS/src/infra/cache"
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

func (b *BaseRepo) returnError(ctx context.Context, err error) *common.Error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common.ErrNotFound(ctx)
	}
	return common.ErrInternal(ctx, err.Error())
}
