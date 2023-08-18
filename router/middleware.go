package router

type Middleware func(Handler) Handler

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) initMiddleware() {
	for reqType, handlers := range r.handlers {
		for i, handler := range handlers {
			for _, m := range r.middlewares {
				r.handlers[reqType][i] = m(handler)
			}
		}
	}
}
