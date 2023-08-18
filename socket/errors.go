package socket

import (
	"context"
	"fmt"
	"log"
)

var (
	ErrPingTimedOut            = fmt.Errorf("ping timed out")
	ErrDisconnectRequest       = fmt.Errorf("disconnected event")
	ErrNilRequest              = fmt.Errorf("nil request")
	ErrUnrecognizedRequestType = fmt.Errorf("unrecognized request type")
)

// handleErrors listens for errors and calls the cancel function when one is received.
func (c *Client) handleErrors(ctx context.Context, cancel context.CancelCauseFunc) {
	defer log.Println("shut down handleErrors")
	log.Println("start handleErrors")
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
