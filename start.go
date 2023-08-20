package socketmode

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

const (
	defaultWaitTimeBetweenRetries = 5 * time.Second
	defaultPingDeadline           = 30 * time.Second
	defaultHandlePingInterval     = 5 * time.Second
)

// Start starts the client.  It contains logic for retries.
func (c *Client) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			if err := c.start(ctx); err != nil {
				// Handle websocket connection error
				if c.retries < c.maxRetries {
					c.retries++
					c.logger.Info("failed to initialize connection", slog.String("error", err.Error()))
					c.logger.Info(fmt.Sprintf("connection retry attempt %d", c.retries))
					time.Sleep(defaultWaitTimeBetweenRetries)
					continue
				}

				return err
			}
		}
	}
}

func (c *Client) start(ctx context.Context) error {
	if err := c.connect(ctx); err != nil {
		return err
	}

	ctx, cancel := context.WithCancelCause(ctx)

	// Listeners
	go c.handleRead(ctx)
	go c.handlePings(ctx)
	go c.handleSend(ctx)
	go c.handleErrors(ctx, cancel)

	c.isConnected = true
	c.logger.Info("client started successfully")

	<-ctx.Done()

	// Close connection
	if err := c.conn.Close(); err != nil {
		return err
	}

	err := context.Cause(ctx)
	c.logger.Info("connection failed with error", slog.Any("error", err))
	return err
}

// connect grabs the `ws` url using the webapi to make the call and opens up a websocket
// connection.
func (c *Client) connect(ctx context.Context) error {
	c.logger.Info("initializing websocket connection")
	_, url, err := c.Api.StartSocketModeContext(ctx)
	if err != nil {
		return err
	}

	if c.debugReconnects {
		url += "&debug_reconnects=true"
	}

	c.conn, _, err = c.dialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.pingTimer = time.NewTimer(c.pingTimeout)
	c.conn.SetPingHandler(pingHandlerFunc(c))
	c.logger.Info("successfully performed websocket handshake")
	return nil
}

// pingHandlerFunc defines the logic of how to handle a ping
func pingHandlerFunc(c *Client) func(string) error {
	return func(h string) error {
		c.retries = 0
		if !c.pingTimer.Stop() {
			<-c.pingTimer.C
		}
		c.pingTimer.Reset(c.pingTimeout)

		return c.conn.WriteControl(websocket.PongMessage, []byte(h), time.Now().Add(defaultPingDeadline))
	}
}

// handlePings is a listener that checks to make sure that our connection is healthy.
func (c *Client) handlePings(ctx context.Context) {
	defer c.logger.Info("shutting down handlePings listener")
	c.logger.Info("starting handlePings listener")
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.pingTimer.C:
			c.errCh <- ErrPingTimeout
		default:
			time.Sleep(defaultHandlePingInterval)
		}
	}
}
