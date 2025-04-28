package request

import "time"

type PaymentCallback struct {
	Token   string  `json:"token"`
	Payment Payment `json:"payment"`
}

type Payment struct {
	TransactionId   string    `json:"transaction_id"`
	Amount          int64     `json:"amount"`
	Content         string    `json:"content"`
	Date            time.Time `json:"date"`
	AccountReceiver string    `json:"account_receiver"`
	Gate            string    `json:"gate"`
}
