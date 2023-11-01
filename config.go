package logwise

import (
	"log/slog"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	// When AddSource is true, the handler adds a ("source", "file:line")
	// attribute to the output indicating the source code position of the log
	// statement. AddSource is false by default to skip the cost of computing
	// this information.
	AddSource bool `cfg:"add_source" mapstructure:"add_source" json:"add_source,omitempty" yaml:"add_source,omitempty"`

	// Level is the minimum enabled logging level.
	Level string `cfg:"level"  mapstructure:"level" json:"level,omitempty" yaml:"level,omitempty"`

	// Encoding sets the logger's encoding. Init values are "json", "text" and "console"
	Encoding string `cfg:"encoding" mapstructure:"encoding" json:"encoding,omitempty" yaml:"encoding,omitempty"`

	// Output is a list of URLs or file paths to write logging output to.
	// See zap.Open for details.
	OutputPaths []string `cfg:"output_paths" mapstructure:"output_paths" json:"output_paths,omitempty" yaml:"output_paths,omitempty"`

	Attributes map[string]any `cfg:"attributes" mapstructure:"attributes" json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

func (cfg *Config) OpenSinks() (zapcore.WriteSyncer, error) {
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = []string{"stderr"}
	}

	sink, _, err := zap.Open(cfg.OutputPaths...)
	return sink, err
}

func (cfg *Config) Opts() *HandlerOptions {
	return &HandlerOptions{HandlerOptions: &slog.HandlerOptions{
		Level:     ToLeveler(cfg.Level),
		AddSource: cfg.AddSource,
	}}
}

func ToLeveler(level string) slog.Leveler {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
