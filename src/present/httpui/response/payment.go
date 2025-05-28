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
