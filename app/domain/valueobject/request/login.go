package request

type Login struct {
	UserName    string `json:"user_name" binding:"required"`    // User Name
	Password    string `json:"password" binding:"required"`     // Password(md5)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // Captcha ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // Captcha Code
}
