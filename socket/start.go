package socket

import (
	"context"
	"log"
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
				log.Printf("retrying attempt %d\n", c.attempts)
				time.Sleep(5 * time.Second)
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
	c.attempts = 0

	ctx, cancel := context.WithCancelCause(ctx)

	// Ping to verify open connection
	go c.checkPing(ctx)

	// Listeners
	go c.listen(ctx)
	go c.handleListen(ctx)
	go c.handleBroadcast(ctx)
	go c.handleErrors(ctx, cancel)

	<-ctx.Done()

	// Clean up
	c.conn.Close()

	return context.Cause(ctx)
}
