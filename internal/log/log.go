// Package log 提供 zap 日志初始化。
// 日志只输出到 stdout,由 systemd 接入 journald(journalctl -u rivalscope 查看),
// 轮转与磁盘限额由 journald 统一管理,应用不再写文件。
package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options 日志初始化参数,来自配置文件 log 段。
type Options struct {
	Level    string // 日志级别:debug/info/warn/error
	Encoding string // console / json(json 更适合 journalctl 检索)
}

// New 根据配置创建 zap.Logger,输出到 stdout。
func New(opt Options) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(opt.Level)); err != nil {
		return nil, fmt.Errorf("解析日志级别失败: %w", err)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoding := opt.Encoding
	if encoding != "json" {
		encoding = "console"
	}
	var encoder zapcore.Encoder
	if encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level), zap.AddCaller()), nil
}
