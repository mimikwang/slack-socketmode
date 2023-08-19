package socket

import "context"

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
	defer c.logger.Debug("shut down handleSend")
	c.logger.Debug("start handleSend")
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
