package models

type LoginWithEmailOTPRequest struct {
	Email  string `json:"email" validate:"required,max=100,email"`
	Source string `json:"source" validate:"required,oneof=web mobile"`
}

type VerifyEmailOTPRequest struct {
	Code  string `json:"code" validate:"required,len=6"`
	Email string `json:"email" validate:"email,required,max=100"`
}

type OIDCLoginRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google"`
	Token    string `json:"token" validate:"required"`
	Source   string `json:"source" validate:"required,oneof=web mobile"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}
