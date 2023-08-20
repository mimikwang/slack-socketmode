package socketmode

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Handler handles incoming events
type Handler func(evt *Event, clt *Client)

// Handle registers generic handlers
func (r *Router) Handle(requestType RequestType, handler Handler, m ...Middleware) {
	handle(r.handlers, requestType, handler, m...)
}

// HandleSlashCommand registers slash command handlers
func (r *Router) HandleSlashCommand(command string, handler Handler, m ...Middleware) {
	handle(r.slashCommandHandlers, command, handler, m...)
}

// HandleEventsApi registers events api handlers
func (r *Router) HandleEventsApi(eventType slackevents.EventsAPIType, handler Handler, m ...Middleware) {
	handle(r.eventsAPIHandlers, eventType, handler, m...)
}

// HandleInteractive registers interactive handlers
func (r *Router) HandleInteractive(interactionType slack.InteractionType, handler Handler, m ...Middleware) {
	handle(r.interactiveHandlers, interactionType, handler, m...)
}

// HandleShortcut registers shortcuts handlers
func (r *Router) HandleShortcut(callbackId string, handler Handler, m ...Middleware) {
	handle(r.shortcutHandlers, callbackId, handler, m...)
}

// HandleBlockAction registers block actions handlers
func (r *Router) HandleBlockAction(actionId string, handler Handler, m ...Middleware) {
	handle(r.blockActionHandlers, actionId, handler, m...)
}

func handle[T comparable](lookup map[T][]Handler, key T, handler Handler, m ...Middleware) {
	if _, found := lookup[key]; !found {
		lookup[key] = []Handler{}
	}
	lookup[key] = append(lookup[key], applyMiddlewares(handler, m...))
}
