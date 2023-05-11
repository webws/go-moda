package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// 新增 options  http.Transport
type ClientOption struct {
	Transport *http.Transport
}

type ClientOptionFunc func(*ClientOption)

func WithClientTransport(transport *http.Transport) ClientOptionFunc {
	return func(option *ClientOption) {
		option.Transport = transport
	}
}

// CallAPI 为 http client 封装,默认使用 otelhttp.NewTransport(http.DefaultTransport)
func CallAPI(ctx context.Context, url string, method string, reqBody interface{}, option ...ClientOptionFunc) ([]byte, error) {
	clientOption := &ClientOption{}
	for _, o := range option {
		o(clientOption)
	}

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	if clientOption.Transport != nil {
		client.Transport = otelhttp.NewTransport(clientOption.Transport)
	}
	var requestBody io.Reader
	if reqBody != nil {
		payload, err := json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}
