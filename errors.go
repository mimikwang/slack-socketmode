package socketmode

import (
	"context"
	"fmt"
)

var (
	ErrPingTimedOut            = fmt.Errorf("ping timed out")
	ErrDisconnectRequest       = fmt.Errorf("disconnected event")
	ErrNilRequest              = fmt.Errorf("nil request")
	ErrUnrecognizedRequestType = fmt.Errorf("unrecognized request type")
	ErrConnClosed              = fmt.Errorf("connection closed")
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
