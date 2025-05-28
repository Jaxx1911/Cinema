package request

import "time"

type PaymentCallback struct {
	Token   string  `json:"token"`
	Payment Payment `json:"payment"`
}

type Payment struct {
	TransactionId   string    `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	Content         string    `json:"content"`
	Date            time.Time `json:"date"`
	AccountReceiver string    `json:"account_receiver"`
	Gate            string    `json:"gate"`
}

type GetPaymentsByCinemaRequest struct {
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02T15:04:05Z"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02T15:04:05Z"`
}

type GetListPaymentRequest struct {
	Page
	Status string `form:"status"`  // success, failed, pending
	SortBy string `form:"sort_by"` // date_desc, date_asc, amount_desc, amount_asc
}
