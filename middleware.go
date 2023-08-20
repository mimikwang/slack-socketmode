package socketmode

// Middleware modifies a handler
type Middleware func(Handler) Handler

// Use registers middlewares
func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) initMiddlewares() {
	applyMiddlewaresToMap(r.middlewares, r.handlers)
	applyMiddlewaresToMap(r.middlewares, r.slashCommandHandlers)
	applyMiddlewaresToMap(r.middlewares, r.eventsAPIHandlers)
	applyMiddlewaresToMap(r.middlewares, r.interactiveHandlers)
	applyMiddlewaresToMap(r.middlewares, r.shortcutHandlers)
	applyMiddlewaresToMap(r.middlewares, r.blockActionHandlers)
}

func applyMiddlewaresToMap[T comparable](middlewares []Middleware, handlersMap map[T][]Handler) {
	for key, handlers := range handlersMap {
		for i, handler := range handlers {
			handlersMap[key][i] = applyMiddlewares(handler, middlewares...)
		}
	}
}

func applyMiddlewares(handler Handler, middlewares ...Middleware) Handler {
	for _, m := range middlewares {
		handler = m(handler)
	}
	return handler
}
