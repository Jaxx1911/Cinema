package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Discount struct {
	Id         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Percentage float64   `json:"percentage"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	UsageLimit int       `json:"usage_limit"`
	IsActive   bool      `json:"is_active"`
}

func ToDiscountResponse(discount domain.Discount) Discount {
	return Discount{
		Id:         discount.ID,
		Code:       discount.Code,
		Percentage: discount.Percentage,
		StartDate:  discount.StartDate.Format("2006-01-02"),
		EndDate:    discount.EndDate.Format("2006-01-02"),
		UsageLimit: discount.UsageLimit,
		IsActive:   discount.IsActive,
	}
}

func ToListDiscountResponse(discount []domain.Discount) []Discount {
	var list []Discount
	for _, d := range discount {
		list = append(list, ToDiscountResponse(d))
	}
	return list
}
