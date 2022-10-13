package mocks

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type MockOkTelegramClient struct{}

func NewMockOkClient() *MockOkTelegramClient {
	return &MockOkTelegramClient{}
}

func (c MockOkTelegramClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "OK", nil
}

type MockFailedTelegramClient struct{}

func NewMockFailedClient() *MockFailedTelegramClient {
	return &MockFailedTelegramClient{}
}

func (c MockFailedTelegramClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "FAIL", errors.New("failed request")
}

type MockHttpClient struct{}

func NewMockHttpClient() *MockHttpClient {
	return &MockHttpClient{}
}

func (c MockHttpClient) PostForm(_ string, _ url.Values) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("OK"))}, nil
}
