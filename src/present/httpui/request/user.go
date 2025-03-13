package request

type UserInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Role  string `json:"role,omitempty"`
}

type Page struct {
	Limit int `json:"limit,omitempty"`
	Page  int `json:"page,omitempty"`
}

func (p *Page) SetDefaults() {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Page <= 0 {
		p.Page = 1
	}
}
