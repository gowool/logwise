package logwise

import (
	"context"
	"log/slog"
	"time"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Discard is a [Handler] which is always disabled and therefore logs nothing.
var Discard slog.Handler = discardHandler{}

var LogLevels = map[zapcore.Level]slog.Level{
	zapcore.DebugLevel:  slog.LevelDebug,
	zapcore.InfoLevel:   slog.LevelInfo,
	zapcore.WarnLevel:   slog.LevelWarn,
	zapcore.ErrorLevel:  slog.LevelError,
	zapcore.DPanicLevel: slog.LevelError,
	zapcore.PanicLevel:  slog.LevelError,
	zapcore.FatalLevel:  slog.LevelError,
}

func NewSlog(logger *zap.Logger) *slog.Logger {
	option := slogzap.Option{
		Level:  LogLevels[logger.Level()],
		Logger: logger.WithOptions(zap.AddCallerSkip(1)),
	}
	return slog.New(option.NewZapHandler())
}

func NewZap(cfg Config) (*zap.Logger, error) {
	cfg.setDefaults()

	var zCfg zap.Config
	switch cfg.Mode {
	case off, none:
		return zap.NewNop(), nil
	case production:
		zCfg = zap.Config{
			Level:       level(cfg.Level),
			Development: false,
			Encoding:    cfg.Encoding,
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      zapcore.OmitKey,
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  zapcore.OmitKey,
				LineEnding:     cfg.LineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     utcEpochTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      cfg.Output,
			ErrorOutputPaths: cfg.ErrorOutput,
		}
	case development:
		zCfg = zap.Config{
			Level:       level(cfg.Level),
			Development: true,
			Encoding:    cfg.Encoding,
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "ts",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      zapcore.OmitKey,
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "msg",
				StacktraceKey:  zapcore.OmitKey,
				LineEnding:     cfg.LineEnding,
				EncodeLevel:    ColoredLevelEncoder,
				EncodeName:     ColoredNameEncoder,
				EncodeTime:     utcISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      cfg.Output,
			ErrorOutputPaths: cfg.ErrorOutput,
		}
	case raw:
		zCfg = zap.Config{
			Level:    level(cfg.Level),
			Encoding: cfg.Encoding,
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey: "message",
				LineEnding: cfg.LineEnding,
			},
			OutputPaths:      cfg.Output,
			ErrorOutputPaths: cfg.ErrorOutput,
		}
	default:
		zCfg = zap.Config{
			Level:    level(cfg.Level),
			Encoding: cfg.Encoding,
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      zapcore.OmitKey,
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  zapcore.OmitKey,
				LineEnding:     cfg.LineEnding,
				EncodeLevel:    ColoredLevelEncoder,
				EncodeName:     ColoredNameEncoder,
				EncodeTime:     utcISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      cfg.Output,
			ErrorOutputPaths: cfg.ErrorOutput,
		}
	}

	if cfg.FileLogger != nil {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FileLogger.Filename,
			MaxSize:    cfg.FileLogger.MaxSize,
			MaxAge:     cfg.FileLogger.MaxAge,
			MaxBackups: cfg.FileLogger.MaxBackups,
			LocalTime:  cfg.FileLogger.LocalTime,
			Compress:   cfg.FileLogger.Compress,
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zCfg.EncoderConfig),
			w,
			zCfg.Level,
		)
		return zap.New(core), nil
	}

	return zCfg.Build()
}

func level(lvl string) zap.AtomicLevel {
	l := zap.NewAtomicLevel()
	_ = l.UnmarshalText([]byte(lvl))
	return l
}

func utcISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05-0700"))
}

func utcEpochTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UTC().UnixNano())
}

type discardHandler struct{}

func (discardHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (discardHandler) Handle(context.Context, slog.Record) error { return nil }
func (d discardHandler) WithAttrs([]slog.Attr) slog.Handler      { return d }
func (d discardHandler) WithGroup(string) slog.Handler           { return d }
