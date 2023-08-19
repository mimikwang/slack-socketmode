package socketmode

// Middleware modifies a handler
type Middleware func(Handler) Handler

// Use registers middlewares
func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) initMiddlewares() {
	applyMiddlewares(r.middlewares, r.handlers)
	applyMiddlewares(r.middlewares, r.slashCommandHandlers)
	applyMiddlewares(r.middlewares, r.eventsAPIHandlers)
}

type mapKey interface {
	string | RequestType
}

func applyMiddlewares[T mapKey](middlewares []Middleware, handlersMap map[T][]Handler) {
	for key, handlers := range handlersMap {
		for i, handler := range handlers {
			for _, m := range middlewares {
				handlersMap[key][i] = m(handler)
			}
		}
	}
}
