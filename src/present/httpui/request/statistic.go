package request

import "time"

type StatisticDateRange struct {
	StartDate time.Time `json:"start_date" form:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" form:"end_date" binding:"required"`
}
