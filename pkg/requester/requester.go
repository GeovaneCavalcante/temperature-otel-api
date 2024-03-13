package requester

import "context"

type Sender interface {
	Send(ctx context.Context, cfg Configuration) (Response, error)
}

type Configuration struct {
	Url         string
	Method      string
	Headers     map[string]string
	Body        interface{}
	ContetType  string
	QueryParams map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
}
