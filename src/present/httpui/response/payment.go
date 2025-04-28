package response

import (
	"TTCS/src/core/domain"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	OrderID       uuid.UUID `json:"order_id"`
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	Amount        int64     `json:"amount"`
	PaymentTime   time.Time `json:"payment_time"`
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
