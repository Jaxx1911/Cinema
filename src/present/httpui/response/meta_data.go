package response

type MetaData struct {
	Data       interface{} `json:"data"`
	TotalCount int64       `json:"total_count"`
}
