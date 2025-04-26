package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
	LocalTime  bool   `json:"local_time"`
}

// An Option configures a Logger.
type Option interface {
	apply(*LogConfig)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*LogConfig)

func (f optionFunc) apply(log *LogConfig) {
	f(log)
}

// AddFileName wraps or replaces the Logger's underlying zapcore.Core.
func FileName(fileName string) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.Filename = fileName
	})
}

func LogLevel(level string) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.Level = level
	})
}

func MaxSize(maxSize int) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.MaxSize = maxSize
	})
}
func MaxAge(maxAge int) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.MaxAge = maxAge
	})
}
func MaxBackups(maxBackups int) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.MaxBackups = maxBackups
	})
}
func Compress(compress bool) Option {
	return optionFunc(func(lc *LogConfig) {
		lc.Compress = compress
	})
}
func (log *LogConfig) WithOptions(opts ...Option) *LogConfig {
	for _, opt := range opts {
		opt.apply(log)
	}
	return log
}

func NewLogConfig(option ...Option) *LogConfig {
	lc := &LogConfig{
		Level:      "info",
		Filename:   "app.log",
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 3,
		Compress:   true,
		LocalTime:  true,
	}
	return lc.WithOptions(option...)
}

func GetZapLog(c *LogConfig) (*zap.Logger, error) {
	// 创建 lumberjack.Logger，配置日志文件的轮转规则
	lumberjackLogger := &lumberjack.Logger{
		Filename:   c.Filename,   // 日志文件路径
		MaxSize:    c.MaxSize,    // 每个日志文件的最大大小（单位：MB）
		MaxBackups: c.MaxBackups, // 最大保留的旧日志文件数量
		MaxAge:     c.MaxAge,     // 旧日志文件的最大保留天数
		Compress:   c.Compress,   // 是否压缩旧日志文件
		LocalTime:  c.LocalTime,
	}

	// 创建 zap 的核心配置
	writerSyncer := zapcore.AddSync(lumberjackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	level, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(encoder, writerSyncer, level)

	// 创建 zap Logger
	logger := zap.New(core, zap.AddCaller()) // 添加调用者信息

	// 记录日志
	logger.Info("This is an info log")
	logger.Error("This is an error log")

	// 确保日志被刷新到文件中
	//logger.Sync()
	return logger, nil
}
