package socket

import (
	"os"

	"golang.org/x/exp/slog"
)

// Opt defines input options for Client
type Opt interface {
	Apply(*Client)
}

// OptDebugReconnects sets the `debugReconnects` flag to true.
type OptDebugReconnects struct{}

func (o OptDebugReconnects) Apply(c *Client) {
	c.debugReconnects = true
}

// OptLogLevel sets the log level
type OptLogLevel struct {
	Level slog.Level
}

func (o OptLogLevel) Apply(c *Client) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: o.Level,
		})
	c.logger = slog.New(handler)
}
