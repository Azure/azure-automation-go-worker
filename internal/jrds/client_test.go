// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package jrds

import (
	"encoding/json"
	"fmt"
	"testing"
)

var (
	baseUri         = "http://jrds.com"
	accountId       = "1d8225ca-97d3-4628-b657-fbb2e0609287"
	sandboxId       = "1d8225ca-97d3-4628-b657-fbb2e0609288"
	workerGroupName = "worker"
)

type BodyMock struct {
	StrProperty string
	IntProperty int
}

var customError = fmt.Errorf("custom error")

type httpClientMock struct {
	get_f  func(url string, headers map[string]string) (responseCode int, body []byte, err error)
	post_f func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
}

func (c httpClientMock) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return c.get_f(url, headers)
}

func (c httpClientMock) Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return c.post_f(url, headers, payload)
}

var getJrdsClient = func(client httpClient) JrdsClient {
	return NewJrdsClient(client, baseUri, accountId, workerGroupName)
}

func TestJrdsClient_issueGetRequest_Returns200(t *testing.T) {
	mock := BodyMock{StrProperty: "string", IntProperty: 2}
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		body, _ = json.Marshal(mock)
		return 200, body, nil
	}}
	client := getJrdsClient(httpClient)

	response := BodyMock{}
	err := client.issueGetRequest(baseUri, &response)
	if err != nil {
		t.Fatalf("unexpected error while calling issueGetRequest")
	}

	if response != mock {
		t.Fatalf("invalid response body")
	}
}

func TestJrdsClient_issueGetRequest_Returns401(t *testing.T) {
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return 401, nil, nil
	}}
	client := getJrdsClient(httpClient)

	err := client.issueGetRequest(baseUri, nil)
	if err == nil {
		t.Fatalf("error is expected for 401 responses")
	}

	switch err.(type) {
	case *RequestAuthorizationError:
		break
	default:
		t.Fatalf("unexpected error type")
	}
}

func TestJrdsClient_issueGetRequest_ReturnsError(t *testing.T) {
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return -1, nil, customError
	}}
	client := getJrdsClient(httpClient)

	err := client.issueGetRequest(baseUri, nil)
	if err == nil {
		t.Fatalf("unexpected error returned by issueGetRequest")
	}
}

func TestJrdsClient_issuePostRequest_Returns200(t *testing.T) {
	mock := BodyMock{StrProperty: "string", IntProperty: 2}
	mock_body := BodyMock{}
	httpClient := httpClientMock{post_f: func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
		json.Unmarshal(payload, &mock_body)
		return 200, nil, nil
	}}
	client := getJrdsClient(httpClient)

	err := client.issuePostRequest(baseUri, mock, nil)
	if err != nil {
		t.Fatalf("unexpected error while calling issuePostRequest")
	}

	if mock_body != mock {
		t.Fatalf("invalid response body")
	}
}

func TestJrdsClient_issuePostRequest_Returns401(t *testing.T) {
	httpClient := httpClientMock{post_f: func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
		return 401, nil, nil
	}}
	client := getJrdsClient(httpClient)

	err := client.issuePostRequest(baseUri, nil, nil)
	if err == nil {
		t.Fatalf("error is expected for 401 responses")
	}

	switch err.(type) {
	case *RequestAuthorizationError:
		break
	default:
		t.Fatalf("unexpected error type")
	}
}

func TestJrdsClient_issuePostRequest_ReturnsError(t *testing.T) {
	httpClient := httpClientMock{post_f: func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
		return -1, nil, customError
	}}
	client := getJrdsClient(httpClient)

	err := client.issuePostRequest(baseUri, nil, nil)
	if err == nil {
		t.Fatalf("unexpected error returned by issuePostRequest")
	}
}

func TestJrdsClient_GetSandboxActions(t *testing.T) {
	mock := SandboxActions{Value: []SandboxAction{{SandboxId: &sandboxId}}}
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		body, _ = json.Marshal(mock)
		return 200, body, nil
	}}
	client := getJrdsClient(httpClient)

	var response SandboxActions
	err := client.GetSandboxActions(&response)
	if err != nil {
		t.Fatalf("unexpected error while calling issueGetRequest")
	}

	if *response.Value[0].SandboxId != sandboxId {
		t.Fatalf("invalid response body")
	}
}
