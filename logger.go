package logwise

import (
	"log/slog"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewSlog(logger zapcore.Core) *slog.Logger {
	return slog.New(zapslog.NewHandler(logger, zapslog.WithCallerSkip(1)))
}

func NewZap(cfg Config) (*zap.Logger, error) {
	cfg.setDefaults()

	if cfg.RollingLogger != nil {
		if err := zap.RegisterSink("rolling", syncFactory(*cfg.RollingLogger)); err != nil {
			return nil, err
		}
	}

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

	return zCfg.Build()
}

type sync struct {
	*lumberjack.Logger
}

func (*sync) Sync() error {
	return nil
}

func syncFactory(cfg RollingLoggerConfig) func(*url.URL) (zap.Sink, error) {
	return func(u *url.URL) (zap.Sink, error) {
		filename := u.Path
		if u.Host == "." || u.Path == "" {
			filename = u.Host + u.Path
		}
		return &sync{&lumberjack.Logger{
			Filename:   filename,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  cfg.LocalTime,
			Compress:   cfg.Compress,
		}}, nil
	}
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
