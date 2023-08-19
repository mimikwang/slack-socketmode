package socket

import (
	"encoding/json"

	"golang.org/x/exp/slog"
)

// Ack sends an acknowledge response to slack.  This can be called concurrently.
func (c *Client) Ack(req *Request, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp := newResponse(req, bytes)
	if err := c.send(resp); err != nil {
		return err
	}
	c.logger.Debug("acknowledged", slog.Any("payload", resp))
	return nil
}
