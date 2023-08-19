package socket

import (
	"context"
)

// Read returns the incoming event.  This is the external API for accessing events.  This function
// is concurrency safe.
func (c *Client) Read() (*Event, error) {
	r := <-c.readCh
	return r.event, r.err
}

type readPackage struct {
	event *Event
	err   error
}

// handleRead is a listener that handles reading events from slack.  Gorilla's read and write
// operations are not concurrency safe, so this should be the only place any read operations are
// called.
//
// More info: https://pkg.go.dev/github.com/gorilla/websocket@v1.4.2#hdr-Concurrency
func (c *Client) handleRead(ctx context.Context) {
	defer c.logger.Debug("shut down handleRead")
	c.logger.Debug("start handleRead")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var req Request
			if err := c.conn.ReadJSON(&req); err != nil {
				// This check is needed here, because when the connection is cut off, it will call
				// `c.conn.ReadJSON` one time, which will result in a read from closed connection
				// error.  That should be the only time that is called when the connection is not
				// started, because the next loop should catch `ctx.Done()`.
				if c.isStarted {
					c.errCh <- err
					c.readCh <- &readPackage{event: nil, err: err}
				}
				return
			}
			evt, err := c.parseRequest(&req)
			c.readCh <- &readPackage{event: evt, err: err}
		}
	}
}

func (c *Client) parseRequest(req *Request) (*Event, error) {
	if req == nil {
		return nil, ErrNilRequest
	}
	if req.Type == RequestTypeDisconnect {
		c.errCh <- ErrDisconnectRequest
	}
	return newEvent(req, context.Background())

}
