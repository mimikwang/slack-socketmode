package socket

import "encoding/json"

// Ack sends an acknowledge response to slack
func (c *Client) Ack(req *Request, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp := newResponse(req, bytes)
	if err := c.send(resp); err != nil {
		return err
	}

	return nil
}
