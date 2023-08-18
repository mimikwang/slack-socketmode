# slack-socketmode
[![Go Reference](https://pkg.go.dev/badge/github.com/mimikwang/slack-socketmode.svg)](https://pkg.go.dev/github.com/mimikwang/slack-socketmode)

Client and Router for Slack's Socketmode

## Introduction

`slack-socketmode` provides a client and router that interact with Slack in [socketmode](https://api.slack.com/apis/connections/socket). It makes use of the community supported slack api client in the [`slack-go/slack` repo](https://github.com/slack-go/slack). This was mostly driven out of the desire to have middleware support for the router.

## Quick Start
Below is a quick example of using the client and router.

```
package main

import (
	"context"
	"fmt"

	"github.com/mimikwang/slack-socketmode/router"
	"github.com/mimikwang/slack-socketmode/socket"
	"github.com/slack-go/slack"
)

const (
	appToken = "YOUR_APP_TOKEN"
	botToken = "YOUR_BOT_TOKEN"
)

func main() {
    // Set up Client
	api := slack.New(botToken, slack.OptionAppLevelToken(appToken))
	client := socket.New(api, socket.OptDebugReconnects{})

    // Set up Router
	r := router.New(client)
	r.Use(dummyMiddleware)
	r.Handle("hello", dummyHandler)

    // Start
	if err := r.Start(context.Background()); err != nil {
		panic(err)
	}
}

func dummyMiddleware(next router.Handler) router.Handler {
	return func(evt *socket.Event, clt *socket.Client) {
		fmt.Println("Stuck in the Middle")
		next(evt, clt)
	}
}

func dummyHandler(evt *socket.Event, clt *socket.Client) {
	fmt.Println("With You")
	clt.Ack(&evt.Request, nil)
}
```