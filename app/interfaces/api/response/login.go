package response

type UserLoginInfo struct {
	UserID   string `json:"user_id"`   // User ID
	UserName string `json:"user_name"` // User Name
	RealName string `json:"real_name"` // Real Name
	Roles    Roles  `json:"roles"`     // Role
}

type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id"` // Captcha ID
}

type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // Access Token
	TokenType   string `json:"token_type"`   // Token Type
	ExpiresAt   int64  `json:"expires_at"`   // Expires At
}
