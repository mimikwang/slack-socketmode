package socketmode

import (
	"os"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

// clientOpt applys options to the client
type clientOpt interface {
	apply(*Client)
}

// ClientOptDialerTimeout sets the timeout duration for the websocket dialer. 0 or negative
// durations are ignored.
type ClientOptDialerTimeout struct {
	Timeout time.Duration
}

func (o ClientOptDialerTimeout) apply(c *Client) {
	if o.Timeout <= 0 {
		return
	}
	c.dialer = &websocket.Dialer{
		HandshakeTimeout: o.Timeout,
	}
}

// ClientOptLogLevel sets the log level. Level should be of type `slog.Level`.
type ClientOptLogLevel struct {
	Level slog.Level
}

func (o ClientOptLogLevel) apply(c *Client) {
	c.logger = slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: o.Level},
		),
	)
}

// ClientOptPingTimeout sets the maximum duration between pings from Slack before timing out. 0 or
// negative durations are ignored.
type ClientOptPingTimeout struct {
	Timeout time.Duration
}

func (o ClientOptPingTimeout) apply(c *Client) {
	if o.Timeout <= 0 {
		return
	}
	c.pingTimeout = o.Timeout
}

// ClientOptDebugReconnects sets the `debugReconnects` flag to true
type ClientOptDebugReconnects struct{}

func (o ClientOptDebugReconnects) apply(c *Client) {
	c.debugReconnects = true
}

// ClientOptMaxRetries sets the number of times to retry connecting to Slack before giving up.
// 0 or negative numbers are ignored.
type ClientOptMaxRetries struct {
	MaxRetries int
}

func (o ClientOptMaxRetries) apply(c *Client) {
	if o.MaxRetries <= 0 {
		return
	}
	c.maxRetries = o.MaxRetries
}

// ClientOptRetryWaitTime sets the duration between retries.  0 or negative durations are ignored.
type ClientOptRetryWaitTime struct {
	WaitTime time.Duration
}

func (o ClientOptRetryWaitTime) apply(c *Client) {
	if o.WaitTime <= 0 {
		return
	}
	c.retryWaitTime = o.WaitTime
}
