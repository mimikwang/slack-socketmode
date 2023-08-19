package socketmode

import (
	"encoding/json"
)

type RequestType string

const (
	RequestTypeDisconnect    = "disconnect"
	RequestTypeHello         = "hello"
	RequestTypeEventsAPI     = "events_api"
	RequestTypeInteractive   = "interactive"
	RequestTypeSlashCommands = "slash_commands"
)

// Request is the incoming requests from slack
type Request struct {
	Type            RequestType     `json:"type"`
	NumConnections  int             `json:"num_connects"`
	DebugInfo       DebugInfo       `json:"debug_info"`
	ConnectionInfo  ConnectionInfo  `json:"connection_info"`
	EnvelopeId      string          `json:"envelope_id"`
	Payload         json.RawMessage `json:"payload"`
	ResponsePayload bool            `json:"response_paload"`
	RetryAttempt    int             `json:"retry_attempt"`
	RetryReason     string          `json:"retry_reason"`
}

type DebugInfo struct {
	Host                      string `json:"host"`
	BuildNumber               int    `json:"build_number"`
	ApproximateConnectionTime int    `json:"approximate_connection_time"`
}

type ConnectionInfo struct {
	AppId string `json:"app_id"`
}
