package socketmode

import (
	"context"
	"fmt"
)

// Custom package errors

var (
	ErrPingTimeout           = fmt.Errorf("timeout waiting for pings from Slack")
	ErrDisconnectRequest     = fmt.Errorf("received disconnected request from Slack")
	ErrNilRequest            = fmt.Errorf("received nil request")
	ErrUnexpectedRequestType = fmt.Errorf("unexpected Request type")
	ErrConnClosed            = fmt.Errorf("websocket connection is closed")
)

// handleErrors listens for errors and calls the cancel function when one is received.
func (c *Client) handleErrors(ctx context.Context, cancel context.CancelCauseFunc) {
	defer c.logger.Info("shutting down handleErrors listener")
	c.logger.Info("starting handleErrors listener")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := <-c.errCh
			cancel(err)
		}
	}
}
