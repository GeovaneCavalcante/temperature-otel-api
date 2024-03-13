package resty

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/GeovaneCavalcante/temperatura-cep/pkg/requester"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Sender struct {
}

func New() *Sender {
	return &Sender{}
}

func (s *Sender) Send(ctx context.Context, cfg requester.Configuration) (requester.Response, error) {

	if _, err := url.ParseRequestURI(cfg.Url); err != nil {
		return requester.Response{}, err
	}

	cli := resty.New().SetHeaders(cfg.Headers).SetHeader("Content-Type", cfg.ContetType).SetTimeout(30 * time.Second)
	cli.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	cli.SetTransport(otelhttp.NewTransport(http.DefaultTransport))
	if cfg.QueryParams != nil {
		cli.SetQueryParams(cfg.QueryParams)
	}

	req := cli.R().SetContext(ctx)

	if cfg.Body != nil {
		req.SetBody(cfg.Body)
	}

	resp, err := req.Execute(strings.ToUpper(cfg.Method), cfg.Url)

	if err != nil {
		return requester.Response{}, err
	}

	response := requester.Response{
		StatusCode: resp.StatusCode(),
		Body:       resp.Body(),
	}

	if resp.IsError() {
		return response, errors.New(string(response.Body))
	}

	return response, nil
}
