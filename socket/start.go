package socket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultWaitTimeBetweenRetries = 5 * time.Second
)

// Start starts the client.  It contains logic for retries.
func (c *Client) Start(ctx context.Context) error {
	for {
		if err := c.start(ctx); err != nil {
			// Check for internal library error types
			switch err {
			case ErrDisconnectRequest:
				// Retry since this is slack's way of telling the client to retry connection
				continue
			}

			// Handle websocket connection error
			if websocket.IsCloseError(err) && c.retries < c.maxRetries {
				c.retries++
				c.logger.Info("connection retry attempt %d\n", c.retries)
				time.Sleep(defaultWaitTimeBetweenRetries)
				continue
			}

			return err
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

	c.isStarted = true
	c.retries = 0

	<-ctx.Done()
	c.isStarted = false // This needs to be set before connection is closed

	// Close connection
	c.conn.Close()
	c.logger.Info("connection failed")

	return context.Cause(ctx)
}
