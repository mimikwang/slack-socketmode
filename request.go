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

// Request is the incoming request from slack
type Request struct {
	Type            RequestType     `json:"type"`
	Reason          string          `json:"reason"`
	NumConnections  int             `json:"num_connections"`
	DebugInfo       DebugInfo       `json:"debug_info"`
	ConnectionInfo  ConnectionInfo  `json:"connection_info"`
	EnvelopeId      string          `json:"envelope_id"`
	Payload         json.RawMessage `json:"payload"`
	ResponsePayload bool            `json:"response_paload"`
	RetryAttempt    int             `json:"retry_attempt"`
	RetryReason     string          `json:"retry_reason"`
}

// DebugInfo is the `debug_info` portion of the incoming Request from Slack
type DebugInfo struct {
	Host                      string `json:"host"`
	BuildNumber               int    `json:"build_number"`
	ApproximateConnectionTime int    `json:"approximate_connection_time"`
}

// ConnectionInfo is the `connection_info` portion of the incoming Request from Slack
type ConnectionInfo struct {
	AppId string `json:"app_id"`
}
