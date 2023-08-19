package socketmode

// Handler handles incoming events
type Handler func(evt *Event, clt *Client)

// Handle registers handlers
func (r *Router) Handle(requestType RequestType, handler Handler) {
	if _, found := r.handlers[requestType]; !found {
		r.handlers[requestType] = []Handler{}
	}
	r.handlers[requestType] = append(r.handlers[requestType], handler)
}

// HandleSlashCommand registers handlers for slash commands
func (r *Router) HandleSlashCommand(command string, handler Handler) {
	if _, found := r.slashCommandHandlers[command]; !found {
		r.slashCommandHandlers[command] = []Handler{}
	}
	r.slashCommandHandlers[command] = append(r.slashCommandHandlers[command], handler)
}

// HandleEventsApi registers handlers for events api
func (r *Router) HandleEventsApi(eventType string, handler Handler) {

}
