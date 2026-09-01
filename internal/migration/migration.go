// Package migration 基于 goose 的嵌入式数据库迁移。
// 迁移 SQL 打包进二进制(go:embed),按数据库方言分目录维护;
// 服务启动时自动执行 Up(幂等),也可通过 `rivalscope migrate` 子命令手动管理。
//
// 职责边界:weibo_account / weibo_feed 由外部爬虫程序创建与维护,
// 本项目只读消费,严禁为其编写迁移;此处只管理 RivalScope 自有表(当前暂无,
// 机制保留待未来微信等模块新增自有表时直接使用)。
package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed all:migrations/postgres
var postgresFS embed.FS

//go:embed all:migrations/mysql
var mysqlFS embed.FS

//go:embed all:migrations/sqlite
var sqliteFS embed.FS

// dialectFS 返回指定驱动对应的迁移文件系统与 goose 方言名。
func dialectFS(driver string) (fs.FS, string, error) {
	switch driver {
	case "postgres":
		sub, err := fs.Sub(postgresFS, "migrations/postgres")
		return sub, "postgres", err
	case "mysql":
		sub, err := fs.Sub(mysqlFS, "migrations/mysql")
		return sub, "mysql", err
	case "sqlite":
		sub, err := fs.Sub(sqliteFS, "migrations/sqlite")
		return sub, "sqlite3", err
	default:
		return nil, "", fmt.Errorf("不支持的数据库驱动: %s", driver)
	}
}

// setup 初始化 goose(绑定方言文件系统与方言)并返回标准库连接。
func setup(db *gorm.DB, dialect string, f fs.FS) (*sql.DB, error) {
	conn, err := db.DB()
	if err != nil {
		return nil, err
	}
	goose.SetBaseFS(f)
	if err := goose.SetDialect(dialect); err != nil {
		return nil, err
	}
	return conn, nil
}

// countMigrations 统计方言目录下的迁移文件数(.sql)。
// 目录可能为空(当前无自有表,机制保留待用),此时迁移操作应视为 no-op。
func countMigrations(f fs.FS) (int, error) {
	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return 0, err
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			n++
		}
	}
	return n, nil
}

// Up 应用全部未执行的迁移(幂等,服务启动时调用);无迁移文件时为 no-op。
func Up(db *gorm.DB, driver string) error {
	f, dialect, err := dialectFS(driver)
	if err != nil {
		return err
	}
	if n, err := countMigrations(f); err != nil || n == 0 {
		return err
	}
	conn, err := setup(db, dialect, f)
	if err != nil {
		return err
	}
	return goose.Up(conn, ".")
}

// Down 回退一个版本(仅手动子命令使用,生产慎用)。
func Down(db *gorm.DB, driver string) error {
	f, dialect, err := dialectFS(driver)
	if err != nil {
		return err
	}
	if n, err := countMigrations(f); err != nil || n == 0 {
		if err != nil {
			return err
		}
		return errors.New("暂无迁移文件,无可回退版本")
	}
	conn, err := setup(db, dialect, f)
	if err != nil {
		return err
	}
	return goose.Down(conn, ".")
}

// Status 输出每个迁移文件的应用状态;无迁移文件时给出提示。
func Status(db *gorm.DB, driver string) error {
	f, dialect, err := dialectFS(driver)
	if err != nil {
		return err
	}
	if n, err := countMigrations(f); err != nil || n == 0 {
		if err != nil {
			return err
		}
		fmt.Println("暂无迁移文件(机制保留,当前无自有表)")
		return nil
	}
	conn, err := setup(db, dialect, f)
	if err != nil {
		return err
	}
	return goose.Status(conn, ".")
}

// Version 返回当前迁移版本号(无版本表/无迁移时为 0)。
func Version(db *gorm.DB, driver string) (int64, error) {
	f, dialect, err := dialectFS(driver)
	if err != nil {
		return 0, err
	}
	if n, err := countMigrations(f); err != nil || n == 0 {
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	conn, err := setup(db, dialect, f)
	if err != nil {
		return 0, err
	}
	version, err := goose.GetDBVersion(conn)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// Run 处理 `rivalscope migrate <cmd>` 子命令。
func Run(db *gorm.DB, driver, cmd string) error {
	switch cmd {
	case "up":
		if err := Up(db, driver); err != nil {
			return err
		}
		fmt.Println("迁移已全部应用")
	case "down":
		if err := Down(db, driver); err != nil {
			return err
		}
		fmt.Println("已回退一个迁移版本")
	case "status":
		return Status(db, driver)
	case "version":
		v, err := Version(db, driver)
		if err != nil {
			return err
		}
		fmt.Printf("当前迁移版本: %d\n", v)
	default:
		return errors.New("不支持的子命令,可用:up / down / status / version")
	}
	return nil
}
