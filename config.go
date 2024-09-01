package logwise

import (
	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	// Mode configures logger based on some default template (development, production, off).
	Mode Mode `json:"mode" yaml:"mode"`

	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling ChannelConfig.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level string `json:"level" yaml:"level"`

	// LineEnding line ending. Default: "\n" for the all modes except production
	LineEnding string `json:"line_ending" yaml:"line_ending"`

	// Encoding sets the logger's encoding. InitDefault values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding" yaml:"encoding"`

	// Output is a list of URLs or file paths to write logging output to.
	// See Open for details.
	Output []string `json:"output" yaml:"output"`

	// ErrorOutput is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutput []string `json:"error_output" yaml:"error_output"`

	// File logger options
	FileLogger *lumberjack.Logger `json:"file_logger,omitempty" yaml:"file_logger,omitempty"`
}

func (cfg *Config) InitDefaults() {
	if cfg.Mode == "" {
		if color.NoColor {
			cfg.Mode = production
		} else {
			cfg.Mode = development
		}
	}
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.LineEnding == "" {
		cfg.LineEnding = zapcore.DefaultLineEnding
	}
	if cfg.Encoding == "" {
		cfg.Encoding = "json"
	}
	if len(cfg.Output) == 0 {
		cfg.Output = []string{"stderr"}
	}
	if len(cfg.ErrorOutput) == 0 {
		cfg.ErrorOutput = []string{"stderr"}
	}
	if cfg.FileLogger != nil {
		if cfg.FileLogger.MaxSize == 0 {
			cfg.FileLogger.MaxSize = 100
		}
		if cfg.FileLogger.MaxAge == 0 {
			cfg.FileLogger.MaxAge = 30
		}
		if cfg.FileLogger.MaxBackups == 0 {
			cfg.FileLogger.MaxBackups = 10
		}
	}
}
