package socket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

// pingHandlerFunc defines the logic of how to handle a ping
func pingHandlerFunc(c *Client) func(string) error {
	return func(h string) error {
		if !c.pingTimer.Stop() {
			<-c.pingTimer.C
		}
		c.pingTimer.Reset(c.maxPingInterval)

		return c.conn.WriteControl(websocket.PongMessage, []byte(h), time.Now().Add(30*time.Second))
	}
}

// handlePings is a listener that checks to make sure that our connection is healthy.
func (c *Client) handlePings(ctx context.Context) {
	defer c.logger.Info("shut down handlePings")
	c.logger.Info("start handlePings")
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.pingTimer.C:
			c.errCh <- ErrPingTimedOut
		default:
			time.Sleep(5 * time.Second)
		}
	}
}
