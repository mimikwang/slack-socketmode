package socketmode

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
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			if err := c.start(ctx); err != nil {
				// Check for internal library error types
				switch err {
				case ErrDisconnectRequest:
					// Retry since this is slack's way of telling the client to retry connection
					time.Sleep(defaultWaitTimeBetweenRetries)
					continue
				}

				// Handle websocket connection error
				if websocket.IsUnexpectedCloseError(err) && c.retries < c.maxRetries {
					c.retries++
					c.logger.Info("connection retry attempt %d\n", c.retries)
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

	c.isStarted = true
	c.retries = 0

	<-ctx.Done()
	c.isStarted = false // This needs to be set before connection is closed

	// Close connection
	c.conn.Close()
	c.logger.Info("connection failed")

	return context.Cause(ctx)
}

// connect grabs the `ws` url using the webapi to make the call and opens up a websocket
// connection.
func (c *Client) connect(ctx context.Context) error {
	c.logger.Info("connecting")
	_, url, err := c.Api.StartSocketModeContext(ctx)
	if err != nil {
		return err
	}

	if c.debugReconnects {
		url += "&debug_reconnects=true"
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.pingTimer = time.NewTimer(c.maxPingInterval)
	c.conn.SetPingHandler(pingHandlerFunc(c))
	c.logger.Info("connected")
	return nil
}
