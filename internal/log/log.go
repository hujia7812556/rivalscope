// Package log 提供 zap 日志初始化,支持 console/json 两种输出与 lumberjack 文件轮转。
package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Options 日志初始化参数,来自配置文件 log 段。
type Options struct {
	Level      string // 日志级别:debug/info/warn/error
	Encoding   string // console / json
	File       string // 日志文件路径,为空则只输出到 stdout
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 旧文件保留天数
	MaxSize    int    // 单文件大小上限(MB)
}

// New 根据配置创建 zap.Logger;同时输出到 stdout 与滚动文件(若配置了 File)。
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

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level),
	}
	// 配置了文件路径时,追加滚动文件输出
	if opt.File != "" {
		lj := &lumberjack.Logger{
			Filename:   opt.File,
			MaxBackups: opt.MaxBackups,
			MaxAge:     opt.MaxAge,
			MaxSize:    opt.MaxSize,
			Compress:   true,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(lj), level))
	}

	return zap.New(zapcore.NewTee(cores...), zap.AddCaller()), nil
}
