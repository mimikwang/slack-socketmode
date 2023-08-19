package socketmode

import (
	"context"
	"fmt"

	"golang.org/x/exp/slog"
)

const (
	defaultParrallel = 5
)

// Router routes incoming requests from slack
type Router struct {
	Client *Client

	parallel int

	middlewares          []Middleware
	handlers             map[RequestType][]Handler
	slashCommandHandlers map[string][]Handler
	eventsAPIHandlers    map[string][]Handler
}

// NewRouter constructs a Router given a socketmode Client
func NewRouter(clt *Client) *Router {
	return &Router{
		Client:               clt,
		parallel:             defaultParrallel,
		middlewares:          []Middleware{},
		handlers:             map[RequestType][]Handler{},
		slashCommandHandlers: map[string][]Handler{},
		eventsAPIHandlers:    map[string][]Handler{},
	}
}

// routerOpt defines options for Router
type routerOpt interface {
	apply(*Router)
}

// OptParralel sets the number of concurrent workers for the router.  This is ignored if the number
// is 0 or negative.
type OptParralel struct {
	Parralel int
}

func (o OptParralel) apply(r *Router) {
	if o.Parralel <= 0 {
		return
	}
	r.parallel = o.Parralel
}

// Start listening to incoming requests
func (r *Router) Start(ctx context.Context) error {
	r.initMiddlewares()
	go r.Client.Start(ctx)

	ctx, cancel := context.WithCancelCause(ctx)
	for i := 0; i < r.parallel; i++ {
		r.Client.logger.Debug(fmt.Sprintf("starting router worker %d", i))
		go func(i int, ctx context.Context, cancel context.CancelCauseFunc) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					evt, err := r.Client.Read()
					if err != nil {
						cancel(err)
					}
					r.Client.logger.Debug(fmt.Sprintf("event routed by worker %d", i))
					r.Client.logger.Debug("route event", slog.Any("event", evt))
					if evt == nil {
						r.Client.logger.Debug("router got nil event")
						return
					}
					r.route(evt)
				}
			}
		}(i, ctx, cancel)
	}

	<-ctx.Done()
	return context.Cause(ctx)
}

// Close cleans up resources
func (r *Router) Close() error {
	return r.Client.Close()
}
