package request

type Discount struct {
	Code       string  `json:"code"`
	Percentage float64 `json:"percentage"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	IsActive   bool    `json:"is_active"`
	UsageLimit int     `json:"usage_limit"` // 0 means unlimited
}
