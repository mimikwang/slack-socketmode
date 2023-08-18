package socket

// Read returns the incoming event.  This is the external API for accessing events.
func (c *Client) Read() (*Event, error) {
	r := <-c.readCh
	return r.event, r.err
}

type readPackage struct {
	event *Event
	err   error
}
