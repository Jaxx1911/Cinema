package request

type UserInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}
