package libhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func do(timeout time.Duration, req *http.Request) (status int, body []byte, err error) {
	defaultHttpClient := http.DefaultClient
	if timeout > 0 {
		defaultHttpClient.Timeout = timeout
	}
	res, err := defaultHttpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	status = res.StatusCode
	return
}

type (
	IHttpRequest interface {
		Request(ctx context.Context, opt RequestOption, out any) (status int, err error)
	}
	HttpClient struct {
	}
)

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (c *HttpClient) Request(ctx context.Context, option RequestOption, out any) (status int, err error) {
	methodStr := methodSelector(option.Method)

	var (
		data io.Reader = nil
		en   Encoder
	)
	if option.Body != nil {
		en = encoder[option.Body.Type]
		data, err = en.Encode(option.Body.Data)
		if err != nil {
			return
		}
	}

	req, err := newRequest(ctx, methodStr, option.Url, data)
	if err != nil {
		err = fmt.Errorf("create request failed: %w", err)
		return
	}

	if option.Header != nil {
		option.Header.Attached(req)
	}

	status, body, err := do(option.RequestTimeout, req)
	if err != nil {
		err = fmt.Errorf("call request failed: %w", err)
		return
	}

	if status >= 200 && status <= 299 {
		if out == nil {
			return
		}
		if err = json.Unmarshal(body, out); err != nil {
			err = fmt.Errorf("parse dst failed: %w", err)
		}
		return
	}
	length := 1024
	if len(body) < length {
		length = len(body) - 1
	}
	return status, fmt.Errorf("call request failed: response=%s", string(body[:length]))
}
