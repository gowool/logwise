package logwise

import (
	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// Mode configures logger based on some default template (development, production, off).
	Mode Mode `json:"mode,omitempty" yaml:"mode,omitempty"`

	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling ChannelConfig.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level string `json:"level,omitempty" yaml:"level,omitempty"`

	// LineEnding line ending. Default: "\n" for the all modes except production
	LineEnding string `json:"lineEnding,omitempty" yaml:"lineEnding,omitempty"`

	// Encoding sets the logger's encoding. InitDefault values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	Encoding string `json:"encoding,omitempty" yaml:"encoding,omitempty"`

	// Output is a list of URLs or file paths to write logging output to.
	// See Open for details.
	Output []string `json:"output,omitempty" yaml:"output,omitempty"`

	// ErrorOutput is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutput []string `json:"errorOutput,omitempty" yaml:"errorOutput,omitempty"`

	// File logger options
	FileLogger *FileLoggerConfig `json:"fileLogger,omitempty" yaml:"fileLogger,omitempty"`
}

func (cfg *Config) setDefaults() {
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
		cfg.FileLogger.setDefaults()
	}
}

type FileLoggerConfig struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `json:"filename,omitempty" yaml:"filename,omitempty"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxAge,omitempty" yaml:"maxAge,omitempty"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxBackups,omitempty" yaml:"maxBackups,omitempty"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localTime,omitempty" yaml:"localTime,omitempty"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `json:"compress,omitempty" yaml:"compress,omitempty"`
}

func (cfg *FileLoggerConfig) setDefaults() {
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 10
	}
}
