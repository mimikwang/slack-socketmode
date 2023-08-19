package socketmode

import (
	"context"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"golang.org/x/exp/slog"
)

// Router routes incoming requests from slack
type Router struct {
	Client *Client

	middlewares          []Middleware
	handlers             map[RequestType][]Handler
	slashCommandHandlers map[string][]Handler
	eventsAPIHandlers    map[string][]Handler
}

// NewRouter constructs a Router given a socketmode Client
func NewRouter(clt *Client) *Router {
	return &Router{
		Client:               clt,
		middlewares:          []Middleware{},
		handlers:             map[RequestType][]Handler{},
		slashCommandHandlers: map[string][]Handler{},
		eventsAPIHandlers:    map[string][]Handler{},
	}
}

// Start listening to incoming requests
func (r *Router) Start(ctx context.Context) error {
	r.initMiddlewares()
	go r.Client.Start(ctx)

	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
			evt, err := r.Client.Read()
			if err != nil {
				return err
			}
			r.Client.logger.Debug("route event", slog.Any("event", evt))
			r.route(evt)
		}
	}
}

// Close cleans up resources
func (r *Router) Close() error {
	return r.Client.Close()
}

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
