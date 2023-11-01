package logwise

import (
	"log"
	"log/slog"
)

// StdLogAdapter can be passed to the http.Server or any place which required standard logger to redirect output
// to the logger plugin
type StdLogAdapter struct {
	log *slog.Logger
}

// Write io.Writer interface implementation
func (s *StdLogAdapter) Write(p []byte) (n int, err error) {
	s.log.Error("internal server error", "error", string(p))
	return len(p), nil
}

// NewStdAdapter constructs StdLogAdapter
func NewStdAdapter(log *slog.Logger) *StdLogAdapter {
	logAdapter := &StdLogAdapter{
		log: log,
	}

	return logAdapter
}

func NewStdLog(base *slog.Logger, prefix string) *log.Logger {
	return log.New(NewStdAdapter(base), prefix, log.Ldate|log.Ltime|log.LUTC)
}
