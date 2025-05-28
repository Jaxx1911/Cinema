package response

import (
	"TTCS/src/core/domain"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID  `json:"id"`
	UserID        *uuid.UUID `json:"user_id"`
	OrderID       *uuid.UUID `json:"order_id"`
	TransactionID string     `json:"transaction_id"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount"`
	PaymentTime   time.Time  `json:"payment_time"`
}

type PaymentDetail struct {
	Date             time.Time `json:"date"`
	UserName         string    `json:"user_name"`
	MovieName        string    `json:"movie_name"`
	RoomName         string    `json:"room_name"`
	Tickets          []string  `json:"tickets"`
	TotalComboAmount float64   `json:"total_combo_amount"`
	TotalAmount      float64   `json:"total_amount"`
	Status           string    `json:"status"`
}

type PaymentWithCustomer struct {
	ID            uuid.UUID  `json:"id"`
	UserID        *uuid.UUID `json:"user_id"`
	OrderID       *uuid.UUID `json:"order_id"`
	TransactionID string     `json:"transaction_id"`
	Status        string     `json:"status"`
	Amount        float64    `json:"amount"`
	PaymentTime   time.Time  `json:"payment_time"`
	Customer      *Customer  `json:"customer,omitempty"`
}

type Customer struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	AvatarUrl string    `json:"avatar_url"`
}

type PaymentCinemaDetail struct {
	ID               uuid.UUID     `json:"id"`
	UserID           *uuid.UUID    `json:"user_id"`
	OrderID          *uuid.UUID    `json:"order_id"`
	TransactionID    string        `json:"transaction_id"`
	Status           string        `json:"status"`
	Amount           float64       `json:"amount"`
	PaymentTime      time.Time     `json:"payment_time"`
	Customer         *Customer     `json:"customer,omitempty"`
	TotalTicketPrice float64       `json:"total_ticket_price"`
	TotalComboPrice  float64       `json:"total_combo_price"`
	Discount         *DiscountInfo `json:"discount,omitempty"`
}

type DiscountInfo struct {
	ID         uuid.UUID `json:"id"`
	Code       string    `json:"code"`
	Percentage float64   `json:"percentage"`
}

type PaymentWithCustomerAndDiscount struct {
	ID               uuid.UUID     `json:"id"`
	UserID           *uuid.UUID    `json:"user_id"`
	OrderID          *uuid.UUID    `json:"order_id"`
	TransactionID    string        `json:"transaction_id"`
	Status           string        `json:"status"`
	Amount           float64       `json:"amount"`
	PaymentTime      time.Time     `json:"payment_time"`
	Customer         *Customer     `json:"customer,omitempty"`
	TotalTicketPrice float64       `json:"total_ticket_price"`
	TotalComboPrice  float64       `json:"total_combo_price"`
	Discount         *DiscountInfo `json:"discount,omitempty"`
}

func ToPaymentResponse(payment domain.Payment) Payment {
	return Payment{
		ID:            payment.ID,
		UserID:        payment.UserID,
		OrderID:       payment.OrderID,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		PaymentTime:   payment.PaymentTime,
	}
}

func ToPaymentsResponse(payments []domain.Payment) []Payment {
	var paymentsResp []Payment
	for _, payment := range payments {
		paymentsResp = append(paymentsResp, ToPaymentResponse(payment))
	}
	return paymentsResp
}

func ToPaymentDetailsResponse(paymentDetails []PaymentDetail) []PaymentDetail {
	return paymentDetails
}

func ToPaymentWithCustomerResponse(payment domain.Payment) PaymentWithCustomer {
	var customer *Customer
	if payment.User.ID != uuid.Nil {
		customer = &Customer{
			ID:        payment.User.ID,
			Name:      payment.User.Name,
			Email:     payment.User.Email,
			Phone:     payment.User.Phone,
			AvatarUrl: payment.User.AvatarUrl,
		}
	}

	return PaymentWithCustomer{
		ID:            payment.ID,
		UserID:        payment.UserID,
		OrderID:       payment.OrderID,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		PaymentTime:   payment.PaymentTime,
		Customer:      customer,
	}
}

func ToPaymentsWithCustomerResponse(payments []domain.Payment) []PaymentWithCustomer {
	var paymentsResp []PaymentWithCustomer
	for _, payment := range payments {
		paymentsResp = append(paymentsResp, ToPaymentWithCustomerResponse(payment))
	}
	return paymentsResp
}

func ToPaymentCinemaDetailResponse(payment domain.Payment, order domain.Order, totalComboPrice float64, discount *DiscountInfo) PaymentCinemaDetail {
	var customer *Customer
	if payment.User.ID != uuid.Nil {
		customer = &Customer{
			ID:        payment.User.ID,
			Name:      payment.User.Name,
			Email:     payment.User.Email,
			Phone:     payment.User.Phone,
			AvatarUrl: payment.User.AvatarUrl,
		}
	}

	// Tính total_ticket_price = total_amount - total_combo_price
	totalTicketPrice := order.TotalPrice - totalComboPrice

	return PaymentCinemaDetail{
		ID:               payment.ID,
		UserID:           payment.UserID,
		OrderID:          payment.OrderID,
		TransactionID:    payment.TransactionID,
		Status:           payment.Status,
		Amount:           payment.Amount,
		PaymentTime:      payment.PaymentTime,
		Customer:         customer,
		TotalTicketPrice: totalTicketPrice,
		TotalComboPrice:  totalComboPrice,
		Discount:         discount,
	}
}

func ToPaymentCinemaDetailsResponse(paymentDetails []PaymentCinemaDetail) []PaymentCinemaDetail {
	return paymentDetails
}

func ToPaymentWithCustomerAndDiscountResponse(payment domain.Payment, order domain.Order, totalComboPrice float64, discount *DiscountInfo) PaymentWithCustomerAndDiscount {
	var customer *Customer
	if payment.User.ID != uuid.Nil {
		customer = &Customer{
			ID:        payment.User.ID,
			Name:      payment.User.Name,
			Email:     payment.User.Email,
			Phone:     payment.User.Phone,
			AvatarUrl: payment.User.AvatarUrl,
		}
	}

	// Tính total_ticket_price = total_amount - total_combo_price
	totalTicketPrice := order.TotalPrice - totalComboPrice

	return PaymentWithCustomerAndDiscount{
		ID:               payment.ID,
		UserID:           payment.UserID,
		OrderID:          payment.OrderID,
		TransactionID:    payment.TransactionID,
		Status:           payment.Status,
		Amount:           payment.Amount,
		PaymentTime:      payment.PaymentTime,
		Customer:         customer,
		TotalTicketPrice: totalTicketPrice,
		TotalComboPrice:  totalComboPrice,
		Discount:         discount,
	}
}

func ToPaymentsWithCustomerAndDiscountResponse(payments []domain.Payment, orderCombos map[uuid.UUID]float64) []PaymentWithCustomerAndDiscount {
	var paymentsResp []PaymentWithCustomerAndDiscount
	for _, payment := range payments {
		if payment.OrderID == nil {
			continue
		}

		// Lấy tổng combo price từ map
		totalComboPrice := orderCombos[*payment.OrderID]

		// Lấy thông tin discount nếu có
		var discountInfo *DiscountInfo
		if payment.Order.DiscountID != nil && payment.Order.Discount != nil {
			discountInfo = &DiscountInfo{
				ID:         payment.Order.Discount.ID,
				Code:       payment.Order.Discount.Code,
				Percentage: payment.Order.Discount.Percentage,
			}
		}

		paymentsResp = append(paymentsResp, ToPaymentWithCustomerAndDiscountResponse(payment, payment.Order, totalComboPrice, discountInfo))
	}
	return paymentsResp
}
