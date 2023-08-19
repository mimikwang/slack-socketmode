package socketmode

import (
	"context"
	"encoding/json"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Event reformats the request into the approrpiate shape and attaches a Context.
type Event struct {
	Request Request
	Data    any
	Context context.Context

	cancel context.CancelFunc
}

// newEvent constructs an Event given a request.
func newEvent(req *Request, ctx context.Context) (*Event, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	var data any
	switch req.Type {
	case RequestTypeHello:
		data = &slack.HelloEvent{}
	case RequestTypeEventsAPI:
		data = &slackevents.EventsAPIEvent{}
	case RequestTypeDisconnect:
		data = &slack.DisconnectedEvent{}
	case RequestTypeSlashCommands:
		data = &slack.SlashCommand{}
	case RequestTypeInteractive:
		data = &slack.InteractionCallback{}
	default:
		return nil, ErrUnrecognizedRequestType
	}

	if len(req.Payload) > 0 {
		if err := json.Unmarshal(req.Payload, data); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Event{Request: *req, Data: data, Context: ctx, cancel: cancel}, nil
}
