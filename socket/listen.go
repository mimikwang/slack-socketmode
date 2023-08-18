package socket

import (
	"context"
	"log"
)

// listen is a listener that listens for events from slack.  Gorilla's read and write operations
// are not concurrency safe, so this should be the only place any read operations are called.
//
// More info: https://pkg.go.dev/github.com/gorilla/websocket@v1.4.2#hdr-Concurrency
func (c *Client) listen(ctx context.Context) {
	defer log.Println("shut down listen")
	log.Println("start listen")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var req Request
			if err := c.conn.ReadJSON(&req); err != nil {
				// Conscious cecision not to send err to errCh here. Since this blocks here and
				// will be called once before the context is cancelled, which will result in
				// killing the reconnect loop.
				//
				// TODO: find a way to handle this more elegantly
				continue
			}
			c.listenCh <- &req
		}
	}
}

// handleListen takes in a Request, reformats it into an Event, and sends it downstream.
func (c *Client) handleListen(ctx context.Context) {
	defer log.Println("shut down handleListen")
	log.Println("start handleListen")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			req := <-c.listenCh

			if req.Type == RequestTypeDisconnect {
				c.errCh <- ErrDisconnectRequest
			}
			evt, err := newEvent(req, context.Background())

			// Send downstream
			c.readCh <- &readPackage{event: evt, err: err}
		}
	}
}
