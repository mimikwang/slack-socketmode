package router

import "context"

func (r *Router) Start(ctx context.Context) error {
	r.initMiddleware()
	go r.Client.Start(ctx)
	for {
		evt, err := r.Client.Read()
		if err != nil {
			return err
		}
		r.routeEvent(evt)
	}
}
