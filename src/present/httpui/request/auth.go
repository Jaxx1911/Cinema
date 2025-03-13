package request

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type GetOTPRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
