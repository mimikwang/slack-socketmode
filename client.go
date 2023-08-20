package socketmode

import (
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack"
	"golang.org/x/exp/slog"
)

const (
	defaultMaxPingInterval = 30 * time.Second
	defaultMaxRetries      = 5
	defaultRetryWaitTime   = 10 * time.Second
)

// Client interacts with Slack in socketmode.
type Client struct {
	// Slack's API Client
	Api *slack.Client

	// Websocket connection and dialer. The isConnected flag signifies if the Client is currently
	// connected to Slack.
	conn        *websocket.Conn
	dialer      *websocket.Dialer
	isConnected bool

	// Structured logger using the `slog` package
	logger *slog.Logger

	// Maximum time to wait between pings from Slack before killing the connection
	pingTimeout time.Duration
	pingTimer   *time.Timer

	// Set to true to debug reconnects. More details at the following:
	// https://api.slack.com/apis/connections/socket#connect
	debugReconnects bool

	// Maximum attempts at retrying connection and wait time between retries
	maxRetries    int
	retryWaitTime time.Duration
	retries       int

	// Communication channels
	readCh chan *readPackage
	sendCh chan *sendPackage
	errCh  chan error
}

// New creates a new socketmode client given a slack api client
func NewClient(api *slack.Client, opts ...clientOpt) *Client {
	c := &Client{
		Api:    api,
		logger: slog.Default(),
		dialer: websocket.DefaultDialer,

		pingTimeout:   defaultMaxPingInterval,
		maxRetries:    defaultMaxRetries,
		retryWaitTime: defaultRetryWaitTime,

		readCh: make(chan *readPackage),
		sendCh: make(chan *sendPackage),
		errCh:  make(chan error),
	}
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
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

// Close cleans up resources
func (c *Client) Close() error {
	if c.isConnected && c.conn != nil {
		// Send a control message to Slack to close the connection
		if err := c.conn.WriteControl(
			websocket.CloseMessage,
			nil,
			time.Now().Add(30*time.Second),
		); err != nil {
			return err
		}
	}
	c.isConnected = false
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
