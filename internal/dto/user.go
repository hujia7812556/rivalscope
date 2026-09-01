package dto

// LoginRequest 登录请求(账号写在配置文件里,不存在注册)。
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginData 登录成功返回的数据。
type LoginData struct {
	AccessToken string `json:"accessToken"`
}

// ProfileData 当前用户信息(来自 JWT claims,不查库)。
type ProfileData struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
