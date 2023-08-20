package socketmode

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// Event contains the original request from Slack as well as the appropriately typed payload.
// It also contains a cancellable context that allows the enforcement that each Event be
// acknowledged at most once.
type Event struct {
	Request Request
	Payload any
	Context context.Context

	cancel context.CancelFunc
}

// newEvent constructs an Event given a Request and Context.
func newEvent(req *Request, ctx context.Context) (*Event, error) {
	if req == nil {
		return nil, ErrNilRequest
	}

	var payload any
	switch req.Type {
	case RequestTypeHello:
		payload = &slack.HelloEvent{}
	case RequestTypeEventsAPI:
		payload = &slackevents.EventsAPICallbackEvent{}
	case RequestTypeDisconnect:
		payload = &slack.DisconnectedEvent{}
	case RequestTypeSlashCommands:
		payload = &slack.SlashCommand{}
	case RequestTypeInteractive:
		payload = &slack.InteractionCallback{}
	default:
		return nil, fmt.Errorf("%s: %s", ErrUnexpectedRequestType, req.Type)
	}

	// Check here as the `hello` payload is empty
	if len(req.Payload) > 0 {
		if err := json.Unmarshal(req.Payload, payload); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Event{Request: *req, Payload: payload, Context: ctx, cancel: cancel}, nil
}
