package router

import (
	"context"
	"time"
)

func (r *Router) Start(ctx context.Context) error {
	r.initMiddleware()
	go r.Client.Start(ctx)

	for {
		if !r.Client.Started() {
			time.Sleep(time.Second)
			continue
		}
		evt, err := r.Client.Read()
		if err != nil {
			return err
		}
		r.routeEvent(evt)
	}
}
