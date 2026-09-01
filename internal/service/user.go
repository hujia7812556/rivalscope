package service

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"rivalscope/internal/auth"
	"rivalscope/internal/config"
	"rivalscope/internal/dto"
)

// UserService 登录业务。
// 账号集合固定(本人与家人),直接来自配置文件 auth.users,无用户表、无注册。
type UserService struct {
	users []config.AuthUser
	jwt   *auth.JWT
}

// NewUserService 创建用户 service;users 为配置写死的账号列表。
func NewUserService(users []config.AuthUser, jwt *auth.JWT) *UserService {
	return &UserService{users: users, jwt: jwt}
}

// isBcrypt 判断密码配置是否为 bcrypt hash($2a$/$2b$/$2y$ 前缀),否则按明文比对。
func isBcrypt(p string) bool {
	return strings.HasPrefix(p, "$2a$") || strings.HasPrefix(p, "$2b$") || strings.HasPrefix(p, "$2y$")
}

// Login 校验用户名密码,成功返回 JWT token(claims 携带身份,无需查库)。
// 用户不存在与密码错误返回同一错误,避免暴露账号是否存在。
func (s *UserService) Login(username, password string) (string, error) {
	for _, u := range s.users {
		if u.Username != username {
			continue
		}
		ok := false
		if isBcrypt(u.Password) {
			ok = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
		} else {
			ok = u.Password == password
		}
		if !ok {
			return "", dto.ErrLoginFailed
		}
		// userId 随机生成仅作会话标识,每次登录不同
		buf := make([]byte, 8)
		if _, err := rand.Read(buf); err != nil {
			return "", dto.ErrInternalServerError
		}
		return s.jwt.GenToken(hex.EncodeToString(buf), u.Username, u.Nickname)
	}
	return "", dto.ErrLoginFailed
}
