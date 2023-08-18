package socket

import (
	"context"

	"golang.org/x/exp/slog"
)

// handleListen is a listener that listens for events from slack.  Gorilla's read and write
// operations are not concurrency safe, so this should be the only place any read operations are
// called.
//
// More info: https://pkg.go.dev/github.com/gorilla/websocket@v1.4.2#hdr-Concurrency
func (c *Client) handleListen(ctx context.Context) {
	defer c.logger.Info("shut down handleListen")
	c.logger.Info("start handleListen")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Avoid reading from closed connection
			if !c.isConnOpened {
				return
			}
			var req Request
			if err := c.conn.ReadJSON(&req); err != nil {
				// Conscious cecision not to send err to errCh here. Since this blocks here and
				// will be called once before the context is cancelled, which will result in
				// killing the reconnect loop.
				//
				// TODO: find a way to handle this more elegantly
				// time.Sleep(time.Second)
				continue
			}
			c.logger.Debug("received request", slog.Any("payload", req))
			c.listenCh <- &req
		}
	}
}
