package handler

import (
	"github.com/gin-gonic/gin"

	"rivalscope/internal/dto"
	"rivalscope/internal/service"
)

// UserHandler 用户接口:登录、个人信息。
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户 handler。
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Login 处理 POST /api/v1/login(账号来自配置文件 auth.users)。
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		dto.Fail(c, err)
		return
	}
	dto.Success(c, dto.LoginData{AccessToken: token})
}

// Profile 处理 GET /api/v1/user,身份直接取自 JWT claims,不查库。
func (h *UserHandler) Profile(c *gin.Context) {
	dto.Success(c, dto.ProfileData{
		Username: c.GetString("username"),
		Nickname: c.GetString("nickname"),
	})
}
