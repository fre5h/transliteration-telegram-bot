package handler

import "errors"

type MockOkHttpClient struct{}

func NewMockOkHttpClient() *MockOkHttpClient {
	return &MockOkHttpClient{}
}

func (c MockOkHttpClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "OK", nil
}

type MockFailedHttpClient struct{}

func NewMockFailedHttpClient() *MockFailedHttpClient {
	return &MockFailedHttpClient{}
}

func (c MockFailedHttpClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "FAIL", errors.New("failed request")
}
