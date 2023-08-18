package socket

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
