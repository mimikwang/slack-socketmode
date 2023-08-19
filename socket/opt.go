package socket

import (
	"os"

	"golang.org/x/exp/slog"
)

// opt defines input options for Client
type opt interface {
	apply(*Client)
}

// OptDebugReconnects sets the `debugReconnects` flag to true.
type OptDebugReconnects struct{}

func (o OptDebugReconnects) apply(c *Client) {
	c.debugReconnects = true
}

// OptLogLevel sets the log level
type OptLogLevel struct {
	Level slog.Level
}

func (o OptLogLevel) apply(c *Client) {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: o.Level,
		})
	c.logger = slog.New(handler)
}

// OptMaxRetries sets the maximum times to retry connection.  0 or negative values will be ignored.
type OptMaxRetries struct {
	MaxRetires int
}

func (o OptMaxRetries) apply(c *Client) {
	if o.MaxRetires <= 0 {
		return
	}
	c.maxRetries = o.MaxRetires
}
