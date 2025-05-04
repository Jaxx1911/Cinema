package request

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
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
