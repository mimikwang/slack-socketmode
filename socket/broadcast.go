package socket

import (
	"context"
	"log"
)

// handleBroadcast is a listener that handles sending events to slack.  Gorilla's read and write
// operations are not concurrency safe, so this should be the only place any write operations are
// called.
//
// More info: https://pkg.go.dev/github.com/gorilla/websocket@v1.4.2#hdr-Concurrency
func (c *Client) handleBroadcast(ctx context.Context) {
	defer log.Println("shut down handleBroadcast")
	log.Println("start handleBroadcast")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			resp := <-c.sendCh
			err := c.conn.WriteJSON(resp.response)
			resp.errCh <- err
		}
	}
}
