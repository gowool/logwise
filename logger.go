package logwise

import "log/slog"

func ToAttrs(data map[string]any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(data))
	for key, value := range data {
		attrs = append(attrs, slog.Any(key, value))
	}
	return attrs
}

func (cfg *Config) Logger(attrs ...slog.Attr) (*slog.Logger, error) {
	syncer, err := cfg.OpenSinks()
	if err != nil {
		return nil, err
	}

	handler := cfg.Opts().NewHandler(syncer, cfg.Encoding)
	handler = handler.WithAttrs(append(ToAttrs(cfg.Attributes), attrs...))

	return slog.New(NewHandlerSyncer(syncer, handler)), nil
}
