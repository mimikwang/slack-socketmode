package socketmode

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (r *Router) route(evt *Event) {
	if r.routeEvent(evt) {
		return
	}
	if r.routeSlashCommand(evt) {
		return
	}
	if r.routeEventsAPI(evt) {
		return
	}
}

func (r *Router) routeEvent(evt *Event) bool {
	if handlers, found := r.handlers[evt.Request.Type]; found {
		applyHandlers(evt, r.Client, handlers)
	}
	return false
}

func (r *Router) routeSlashCommand(evt *Event) bool {
	if evt.Request.Type != RequestTypeSlashCommands {
		return false
	}

	data, ok := evt.Data.(*slack.SlashCommand)
	if !ok {
		return false
	}

	if handlers, found := r.slashCommandHandlers[data.Command]; found {
		applyHandlers(evt, r.Client, handlers)
	}
	return false
}

func (r *Router) routeEventsAPI(evt *Event) bool {
	if evt.Request.Type != RequestTypeEventsAPI {
		return false
	}

	data, ok := evt.Data.(*slackevents.EventsAPIEvent)
	if !ok {
		return false
	}

	if handlers, found := r.eventsAPIHandlers[data.InnerEvent.Type]; found {
		applyHandlers(evt, r.Client, handlers)
	}
	return false
}

func applyHandlers(evt *Event, clt *Client, handlers []Handler) bool {
	for _, handler := range handlers {
		select {
		case <-evt.Context.Done():
			return true
		default:
			handler(evt, clt)
		}
	}
	return false
}
