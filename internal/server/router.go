// Package server 组装路由与依赖(项目规模小,手动组装替代依赖注入框架)。
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"rivalscope/internal/auth"
	"rivalscope/internal/config"
	"rivalscope/internal/handler"
	"rivalscope/internal/middleware"
	"rivalscope/internal/repo"
	"rivalscope/internal/service"
)

// New 创建 gin 引擎并注册全部 API 路由。
func New(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 依赖组装:repo -> service -> handler(登录账号来自配置,无用户表)
	jwtAuth := auth.NewJWT(cfg.Security.JWT.Key, cfg.Security.JWT.ExpireHours)
	accountRepo := repo.NewWeiboAccountRepo(db)
	msgRepo := repo.NewWeiboMsgRepo(db)
	statsRepo := repo.NewWeiboStatisticsRepo(db)

	userHandler := handler.NewUserHandler(service.NewUserService(cfg.Auth.Users, jwtAuth))
	accountHandler := handler.NewWeiboAccountHandler(service.NewWeiboAccountService(accountRepo))
	msgHandler := handler.NewWeiboMsgHandler(service.NewWeiboMsgService(msgRepo))
	statsHandler := handler.NewWeiboStatisticsHandler(service.NewWeiboStatisticsService(statsRepo))

	r := gin.New()
	r.Use(gin.Recovery(), requestLogger(logger))

	api := r.Group("/api/v1")
	// 公开接口:登录
	api.POST("/login", userHandler.Login)

	// 以下接口全部需要 JWT
	authed := api.Group("", middleware.JWT(jwtAuth))
	authed.GET("/user", userHandler.Profile)
	authed.GET("/weibo/accounts", accountHandler.GetList)
	authed.POST("/weibo/accounts", accountHandler.Create)
	authed.PATCH("/weibo/accounts/:id", accountHandler.Update)
	authed.DELETE("/weibo/accounts/:id", accountHandler.Delete)
	authed.GET("/weibo/messages", msgHandler.GetList)
	authed.GET("/weibo/statistics", statsHandler.GetList)

	return r
}

// requestLogger 用 zap 记录每个请求的方法、路径、状态码与耗时。
func requestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// Run 启动 HTTP 服务并支持优雅退出(ctx 取消时最多等待 5 秒)。
func Run(ctx context.Context, cfg *config.Config, logger *zap.Logger, r *gin.Engine) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("HTTP 服务启动", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
