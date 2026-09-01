// RivalScope(竞品观察)—— 查看竞品社媒账号数据的单体服务。
// React 前端构建产物通过 go:embed 嵌入,单二进制 + 单配置文件即可部署。
// 登录账号直接写在配置 auth.users(用户量固定),无用户表。
//
// 用法:
//
//	rivalscope -config config/config.yaml                       启动 HTTP 服务(启动时自动执行迁移)
//	rivalscope migrate up|down|status|version -config ...       手动管理数据库迁移
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"rivalscope/internal/config"
	"rivalscope/internal/log"
	"rivalscope/internal/migration"
	"rivalscope/internal/repo"
	"rivalscope/internal/server"
)

func main() {
	// 子命令:migrate 手动管理数据库迁移
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := runMigrate(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "迁移操作失败:", err)
			os.Exit(1)
		}
		return
	}

	confPath := flag.String("config", "config/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*confPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger, err := log.New(log.Options{
		Level:    cfg.Log.Level,
		Encoding: cfg.Log.Encoding,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer logger.Sync()

	db, err := repo.NewDB(cfg.Data.DB.Driver, cfg.Data.DB.DSN)
	if err != nil {
		logger.Fatal("数据库初始化失败", zap.Error(err))
	}
	// 启动时自动应用迁移(幂等;只管理 RivalScope 自有表,不触碰爬虫维护的微博表)
	if err := migration.Up(db, cfg.Data.DB.Driver); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	r := server.New(cfg, db, logger)
	registerWeb(r)

	// 监听退出信号,优雅停机
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, cfg, logger, r); err != nil {
		logger.Fatal("HTTP 服务异常退出", zap.Error(err))
	}
	logger.Info("HTTP 服务已停止")
}

// runMigrate 处理 migrate 子命令:rivalscope migrate <up|down|status|version>。
func runMigrate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定子命令:up / down / status / version")
	}
	sub := args[0]

	fs := flag.NewFlagSet("migrate "+sub, flag.ExitOnError)
	confPath := fs.String("config", "config/config.yaml", "配置文件路径")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	cfg, err := config.Load(*confPath)
	if err != nil {
		return err
	}
	db, err := repo.NewDB(cfg.Data.DB.Driver, cfg.Data.DB.DSN)
	if err != nil {
		return err
	}
	return migration.Run(db, cfg.Data.DB.Driver, sub)
}
