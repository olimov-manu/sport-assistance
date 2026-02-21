package responses

type SendOTPResponse struct {
	OTPSent bool   `json:"otp_sent"`
	Message string `json:"message"`
}

type ConfirmOTPResponse struct {
	OTPConfirmed bool   `json:"otp_confirmed"`
	IsRegistered bool   `json:"is_registered"`
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
