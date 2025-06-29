package resetpassword

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyCodeRequest struct {
	Email     string `json:"email"`
	ResetCode string `json:"code"`
}

type UpdatePasswordRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type UpdatePasswordResponse struct {
	Email string `json:"email"`
}
