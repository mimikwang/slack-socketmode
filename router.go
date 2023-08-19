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
	slog.Info("starting router")
	r.initMiddlewares()
	routerCtx, cancel := context.WithCancelCause(ctx)

	go func(cancel context.CancelCauseFunc) {
		if err := r.Client.Start(ctx); err != nil {
			cancel(err)
		}
	}(cancel)

	for i := 0; i < r.parallel; i++ {
		go r.startWorker(i, routerCtx)
	}

	slog.Info("router started successfully")
	<-routerCtx.Done()

	err := context.Cause(routerCtx)
	slog.Error("router stopped with error", slog.String("error", err.Error()))
	return err
}

func (r *Router) startWorker(i int, ctx context.Context) {
	defer r.Client.logger.Info(fmt.Sprintf("shutting down router worker %d", i))
	r.Client.logger.Info(fmt.Sprintf("starting router worker %d", i))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			evt, err := r.Client.Read()
			if err != nil {
				r.Client.logger.Debug("client read error", slog.Any("error", err))
				continue
			}

			r.Client.logger.Info(fmt.Sprintf("event routed by worker %d", i))
			if evt == nil {
				r.Client.logger.Info("router got nil event")
				continue
			}
			r.Client.logger.Debug(
				"route event",
				slog.Any("event_type", evt.Request.Type),
				slog.Any("event_data", evt.Data),
			)
			r.route(evt)
		}
	}
}

// Close cleans up resources
func (r *Router) Close() error {
	return r.Client.Close()
}
