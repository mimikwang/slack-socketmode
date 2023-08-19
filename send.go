package socketmode

import (
	"context"
	"encoding/json"

	"golang.org/x/exp/slog"
)

// Ack sends an acknowledge response to slack.  This can be called concurrently.
func (c *Client) Ack(evt *Event, payload any) error {
	if evt == nil {
		return ErrNilRequest
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp := newResponse(&evt.Request, bytes)
	if err := c.send(resp); err != nil {
		return err
	}

	// Each event should only be acknowledged once
	evt.cancel()
	c.logger.Debug(
		"acknowledged",
		slog.Any("envelope_id", resp.EnvelopeId),
		slog.Any("payload", string(resp.Payload)),
	)
	return nil
}

// Send sends a response to slack. This is not exposed to the user since all responses should
// be called via `Ack`.
func (c *Client) send(resp *Response) error {
	errCh := make(chan error)
	c.sendCh <- &sendPackage{response: resp, errCh: errCh}
	err := <-errCh
	close(errCh)

	return err
}

type sendPackage struct {
	response *Response
	errCh    chan error
}

// handleSend is a listener that handles sending events to slack.  Gorilla's read and write
// operations are not concurrency safe, so this should be the only place any write operations are
// called.
//
// More info: https://pkg.go.dev/github.com/gorilla/websocket@v1.4.2#hdr-Concurrency
func (c *Client) handleSend(ctx context.Context) {
	defer c.logger.Info("shutting down handleSend listener")
	c.logger.Info("starting handleSend listener")
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
