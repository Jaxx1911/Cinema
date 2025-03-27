package request

type GetCinemaRequest struct {
	City string `json:"city" binding:"required" form:"city"`
}

func (g *GetCinemaRequest) MappingCity() {
	switch g.City {
	case "hanoi":
		g.City = "Hà Nội"
	case "hcm":
		g.City = "HCM"
	case "danang":
		g.City = "Đà Nẵng"
	}
}
