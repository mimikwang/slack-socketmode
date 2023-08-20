package socketmode

import (
	"encoding/json"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func (r *Router) route(evt *Event) {
	switch evt.Request.Type {
	case RequestTypeEventsAPI:
		if r.routeEventsAPI(evt) {
			return
		}
	case RequestTypeSlashCommands:
		if r.routeSlashCommand(evt) {
			return
		}
	case RequestTypeInteractive:
		if r.routeInteractive(evt) {
			return
		}
	}
	r.routeEvent(evt)
}

func (r *Router) routeEvent(evt *Event) bool {
	return applyHandlersMap(evt, r.client, r.handlers, evt.Request.Type)
}

func (r *Router) routeSlashCommand(evt *Event) bool {
	data := evt.Payload.(*slack.SlashCommand)
	return applyHandlersMap(evt, r.client, r.slashCommandHandlers, data.Command)
}

type eventType struct {
	Type string `json:"type"`
}

func (r *Router) routeEventsAPI(evt *Event) bool {
	payload := evt.Payload.(*slackevents.EventsAPICallbackEvent)

	var innerType eventType
	if err := json.Unmarshal(*payload.InnerEvent, &innerType); err != nil {
		return false
	}

	return applyHandlersMap(evt, r.client, r.eventsAPIHandlers, slackevents.EventsAPIType(innerType.Type))
}

func (r *Router) routeInteractive(evt *Event) bool {
	data := evt.Payload.(*slack.InteractionCallback)

	switch data.Type {
	case slack.InteractionTypeShortcut, slack.InteractionTypeMessageAction:
		if applyHandlersMap(evt, r.client, r.shortcutHandlers, data.CallbackID) {
			return true
		}
	case slack.InteractionTypeBlockActions:
		for _, action := range data.ActionCallback.BlockActions {
			if applyHandlersMap(evt, r.client, r.blockActionHandlers, action.ActionID) {
				return true
			}
		}
	}

	return applyHandlersMap(evt, r.client, r.interactiveHandlers, data.Type)
}

func applyHandlersMap[T comparable](evt *Event, clt *Client, lookup map[T][]Handler, key T) bool {
	if handlers, found := lookup[key]; found {
		return applyHandlers(evt, clt, handlers)
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
