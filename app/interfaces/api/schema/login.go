package schema

type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`    // User Name
	Password    string `json:"password" binding:"required"`     // Password(md5)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // Captcha ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // Captcha Code
}

type UserLoginInfo struct {
	UserID   string `json:"user_id"`   // User ID
	UserName string `json:"user_name"` // User Name
	RealName string `json:"real_name"` // Real Name
	Roles    Roles  `json:"roles"`     // Role
}

type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // Old Password
	NewPassword string `json:"new_password" binding:"required"` // New Password
}

type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id"` // Captcha ID
}

type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // Access Token
	TokenType   string `json:"token_type"`   // Token Type
	ExpiresAt   int64  `json:"expires_at"`   // Expires At
}
