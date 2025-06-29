package resetpassword

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyCodeRequest struct {
	Email     string `json:"email"`
	ResetCode string `json:"code"`
}
