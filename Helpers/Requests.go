package Helpers

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

var client = &fasthttp.Client{
	MaxConnDuration: time.Duration(Timeout) * time.Millisecond,
}

type RequestOptions struct {
	Headers        map[string]string
	CustomizeFunc  func(*fasthttp.Request)
	LoggingEnabled bool
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func SendRequest(url, method, payload string, options RequestOptions) (*Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()

	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetUserAgent("Mozilla/5.0")

	if options.CustomizeFunc != nil {
		options.CustomizeFunc(req)
	}

	for key, value := range options.Headers {
		req.Header.Add(key, value)
	}

	if payload != "" {
		req.SetBodyString(payload)
	}

	err := client.Do(req, resp)
	if err != nil {
		if options.LoggingEnabled {
			fmt.Printf("Request failed for URL %s with error: %v\n", url, err)
		}
		return nil, err
	}

	statusCode := resp.StatusCode()
	if statusCode >= fasthttp.StatusBadRequest {
		if options.LoggingEnabled {
			fmt.Printf("HTTP request failed for URL %s with status code %d\n", url, statusCode)
		}
		return nil, fmt.Errorf("HTTP request failed with status code %d", statusCode)
	}

	response := &Response{
		StatusCode: statusCode,
		Headers:    make(map[string]string),
	}

	resp.Header.VisitAll(func(key, value []byte) {
		response.Headers[string(key)] = string(value)
	})

	response.Body = resp.Body()

	if options.LoggingEnabled {
		//fmt.Printf("Request successful for URL %s with status code %d\n", url, statusCode)
	}

	return response, nil
}
