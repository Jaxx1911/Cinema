package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type Discount struct {
	Id         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Percentage float64   `json:"percentage"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

func ToDiscountResponse(discount domain.Discount) Discount {
	return Discount{
		Id:         discount.ID,
		Code:       discount.Code,
		Percentage: discount.Percentage,
		StartDate:  discount.StartDate,
		EndDate:    discount.EndDate,
	}
}

func ToListDiscountResponse(discount []domain.Discount) []Discount {
	var list []Discount
	for _, d := range discount {
		list = append(list, ToDiscountResponse(d))
	}
	return list
}
