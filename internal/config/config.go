// Package config 负责 RivalScope 的配置加载与校验。
package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置结构,与 config/*.yaml 字段一一对应。
type Config struct {
	Env  string `mapstructure:"env"`
	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"http"`

	Security struct {
		JWT struct {
			Key         string `mapstructure:"key"`
			ExpireHours int    `mapstructure:"expire_hours"`
		} `mapstructure:"jwt"`
	} `mapstructure:"security"`

	// Auth 登录账号(用户量固定,直接写在配置里,不建用户表)。
	Auth struct {
		Users []AuthUser `mapstructure:"users"`
	} `mapstructure:"auth"`

	Data struct {
		DB struct {
			Driver string `mapstructure:"driver"` // postgres / mysql / sqlite
			DSN    string `mapstructure:"dsn"`
		} `mapstructure:"db"`
	} `mapstructure:"data"`

	Log struct {
		Level    string `mapstructure:"level"`    // 日志级别:debug/info/warn/error
		Encoding string `mapstructure:"encoding"` // console / json,输出到 stdout 由 journald 接管
	} `mapstructure:"log"`
}

// AuthUser 配置写死的登录账号。
type AuthUser struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"` // 明文,或 $2a$/$2b$/$2y$ 开头的 bcrypt hash
	Nickname string `mapstructure:"nickname"`
}

// Load 从指定路径读取 yaml 配置并做基础校验与默认值兜底。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	if strings.TrimSpace(c.Security.JWT.Key) == "" {
		return nil, fmt.Errorf("配置项 security.jwt.key 不能为空")
	}
	if c.Security.JWT.ExpireHours <= 0 {
		c.Security.JWT.ExpireHours = 24 * 90 // 默认 90 天
	}
	if c.HTTP.Port <= 0 {
		c.HTTP.Port = 17317
	}
	if c.HTTP.Host == "" {
		c.HTTP.Host = "127.0.0.1"
	}
	return &c, nil
}
