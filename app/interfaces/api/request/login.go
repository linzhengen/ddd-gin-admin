package request

type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`    // User Name
	Password    string `json:"password" binding:"required"`     // Password(md5)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // Captcha ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // Captcha Code
}

type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // Old Password
	NewPassword string `json:"new_password" binding:"required"` // New Password
}
