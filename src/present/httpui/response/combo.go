package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
)

type Combo struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BannerUrl   string    `json:"banner_url"`
	Price       int64     `json:"price"`
}

type ComboWithQuantity struct {
	Combo
	Quantity int64 `json:"quantity"`
}

func ToComboResponse(combo *domain.Combo) Combo {
	return Combo{
		Id:          combo.ID,
		Name:        combo.Name,
		Description: combo.Description,
		BannerUrl:   combo.BannerUrl,
		Price:       int64(combo.Price),
	}
}

func ToListComboResponse(comboList []*domain.Combo) []Combo {
	var listCombo []Combo
	for _, combo := range comboList {
		listCombo = append(listCombo, ToComboResponse(combo))
	}
	return listCombo
}

func ToListComboWithQuantity(comboList []domain.OrderCombo) []ComboWithQuantity {
	var listComboWithQuantity []ComboWithQuantity
	for _, combo := range comboList {
		listComboWithQuantity = append(listComboWithQuantity, ComboWithQuantity{
			Combo:    ToComboResponse(&combo.Combo),
			Quantity: int64(combo.Quantity),
		})
	}
	return listComboWithQuantity
}
