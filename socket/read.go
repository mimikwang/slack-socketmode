package socket

import (
	"context"

	"golang.org/x/exp/slog"
)

// Read returns the incoming event.  This is the external API for accessing events.
func (c *Client) Read() (*Event, error) {
	r := <-c.readCh
	return r.event, r.err
}

type readPackage struct {
	event *Event
	err   error
}

// handleRead takes in a Request, reformats it into an Event, and sends it downstream.
func (c *Client) handleRead(ctx context.Context) {
	defer c.logger.Info("shut down handleRead")
	c.logger.Info("start handleRead")
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

			c.logger.Debug("new read event", slog.Any("payload", evt))

			// Send downstream
			c.readCh <- &readPackage{event: evt, err: err}
		}
	}
}
