package socket

// Close cleans up resources
func (c *Client) Close() error {
	return c.conn.Close()
}
