package transaction

import (
	"TTCS/src/common/adapters"
	"TTCS/src/infra/cache"
	"TTCS/src/infra/repo"
	"gorm.io/gorm"
)

type Provider struct {
	db    *gorm.DB
	cache *cache.RedisCache
}

func NewTransactionProvider(db *gorm.DB, redis *cache.RedisCache) *Provider {
	return &Provider{
		cache: redis,
		db:    db,
	}
}

func (p *Provider) Transact(txFunc func(adapter adapters.Adapters) error) error {
	err := runInTx(p.db, func(tx *gorm.DB) error {
		baseRepo := repo.NewBaseRepo(tx, p.cache)
		adapter := adapters.Adapters{
			OrderRepo:    repo.NewOrderRepo(baseRepo),
			ShowtimeRepo: repo.NewShowtimeRepo(baseRepo),
			CinemaRepo:   repo.NewCinemaRepo(baseRepo),
			ComboRepo:    repo.NewComboRepo(baseRepo),
			DiscountRepo: repo.NewDiscountRepo(baseRepo),
			GenreRepo:    repo.NewGenreRepo(baseRepo),
			MovieRepo:    repo.NewMovieRepo(baseRepo),
			Payment:      repo.NewPaymentRepo(baseRepo),
			RoomRepo:     repo.NewRoomRepo(baseRepo),
			SeatRepo:     repo.NewSeatRepo(baseRepo),
			TicketRepo:   repo.NewTicketRepo(baseRepo),
			UserRepo:     repo.NewUserRepo(baseRepo),
		}

		return txFunc(adapter)
	})
	return err
}

func runInTx(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()

	err := fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
