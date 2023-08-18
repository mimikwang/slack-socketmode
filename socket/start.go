package socket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
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
			if websocket.IsCloseError(err) && c.attempts < c.maxAttempts {
				c.attempts++
				c.logger.Info("connection retry attempt %d\n", c.attempts)
				time.Sleep(5 * time.Second)
				continue
			}

			return err
		}
	}
}

// Started returns true of the client has started
func (c *Client) Started() bool {
	return c.isStarted
}

func (c *Client) start(ctx context.Context) error {
	if err := c.connect(ctx); err != nil {
		return err
	}

	ctx, cancel := context.WithCancelCause(ctx)

	// Listeners
	go c.handlePings(ctx)
	go c.handleSend(ctx)
	go c.handleErrors(ctx, cancel)

	c.isStarted = true
	c.attempts = 0

	<-ctx.Done()

	// Clean up
	c.conn.Close()
	c.isStarted = false
	c.logger.Info("connection failed")

	return context.Cause(ctx)
}
