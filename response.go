package socketmode

import "encoding/json"

// Response is what slack expects as a response. It should be called with `Ack` once.
type Response struct {
	EnvelopeId string          `json:"envelope_id"`
	Payload    json.RawMessage `json:"payload"`
}

func newResponse(req *Request, payload json.RawMessage) *Response {
	return &Response{
		EnvelopeId: req.EnvelopeId,
		Payload:    payload,
	}
}
