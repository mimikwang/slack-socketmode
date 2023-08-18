package socket

// Opt defines input options for Client
type Opt interface {
	Apply(*Client)
}

// OptDebugReconnects sets the `debugReconnects` flag to true.
type OptDebugReconnects struct{}

func (o OptDebugReconnects) Apply(c *Client) {
	c.debugReconnects = true
}
