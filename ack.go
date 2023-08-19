package socketmode

import (
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
	c.logger.Debug("acknowledged", slog.Any("payload", resp))
	return nil
}
