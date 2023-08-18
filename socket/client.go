package socket

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/slack-go/slack"
)

const (
	defaultMaxPingInterval = 30 * time.Second
	defaultMaxAttempts     = 5
)

// Client interacts with slack in socketmode
type Client struct {
	// Slack's API Client
	Api *slack.Client

	// Websocket client
	conn *websocket.Conn

	// Max time between pings before timing out
	maxPingInterval time.Duration
	pingTimer       *time.Timer

	// Set to true to debug reconnects.
	// More details here: https://api.slack.com/apis/connections/socket#connect
	debugReconnects bool

	// Maximum attempts at reconnecting
	maxAttempts int
	attempts    int

	readCh   chan *readPackage
	listenCh chan *Request
	sendCh   chan *sendPackage
	errCh    chan error
}

// New creates a new socketmode client given a slack api client
func New(api *slack.Client, opts ...Opt) *Client {
	c := &Client{
		Api:             api,
		maxPingInterval: defaultMaxPingInterval,
		maxAttempts:     defaultMaxAttempts,

		readCh:   make(chan *readPackage),
		listenCh: make(chan *Request),
		sendCh:   make(chan *sendPackage),
		errCh:    make(chan error),
	}
	for _, opt := range opts {
		opt.Apply(c)
	}
	return c
}
