package socket

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

// connect grabs the `ws` url using the webapi to make the call and opens up a websocket
// connection.
func (c *Client) connect(ctx context.Context) error {
	c.logger.Info("connecting")
	_, url, err := c.Api.StartSocketModeContext(ctx)
	if err != nil {
		return err
	}

	if c.debugReconnects {
		url += "&debug_reconnects=true"
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.pingTimer = time.NewTimer(c.maxPingInterval)
	c.conn.SetPingHandler(pingHandlerFunc(c))
	c.logger.Info("connected")
	return nil
}
