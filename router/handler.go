package router

import "github.com/mimikwang/slack-socketmode/socket"

type Handler func(evt *socket.Event, clt *socket.Client)

func (r *Router) Handle(requestType socket.RequestType, handler Handler) {
	if _, found := r.handlers[requestType]; !found {
		r.handlers[requestType] = []Handler{}
	}
	r.handlers[requestType] = append(r.handlers[requestType], handler)
}
