package socket

import (
	"context"
	"log"
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

// checkPing is a listener that checks to make sure that our connection is healthy.
func (c *Client) checkPing(ctx context.Context) {
	defer log.Println("shut down checkPing")
	log.Println("start checkPing")
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
