package router

import (
	"github.com/mimikwang/slack-socketmode/socket"
)

type Router struct {
	Client *socket.Client

	middlewares []Middleware
	handlers    map[socket.RequestType][]Handler
}

func New(clt *socket.Client) *Router {
	return &Router{
		Client:      clt,
		middlewares: []Middleware{},
		handlers:    map[socket.RequestType][]Handler{},
	}
}

func (r *Router) routeEvent(evt *socket.Event) {
	if handlers, found := r.handlers[evt.Request.Type]; found {
		for _, handler := range handlers {
			go handler(evt, r.Client)
		}
	}
}
