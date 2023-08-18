package socket

import (
	"context"
)

// Read returns the incoming event.  This is the external API for accessing events.
func (c *Client) Read() (*Event, error) {
	req, err := c.readRequest()
	if err != nil {
		return nil, err
	}
	return c.parseRequest(req)
}

type readPackage struct {
	event *Event
	err   error
}

func (c *Client) readRequest() (*Request, error) {
	var req Request
	if !c.Started() {
		return nil, ErrConnClosed
	}
	if err := c.conn.ReadJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
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
