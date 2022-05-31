package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fre5h/transliteration-telegram-bot/internal/model"
)

func TestParseTelegramRequest(t *testing.T) {
	var update = model.Update{
		UpdateId: 1,
		Message: model.Message{
			Text: "привіт",
			Chat: model.Chat{
				Id: 1,
			},
		},
	}

	requestBody, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Failed to marshal update in json, got %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))

	var updateToTest, errParse = parseTelegramRequest(req)
	if errParse != nil {
		t.Errorf("Expected a <nil> error, got %s", errParse.Error())
	}

	if *updateToTest != update {
		t.Errorf("Expected update %v, got %v", update, updateToTest)
	}
}

func TestParseTelegramRequestInvalid(t *testing.T) {
	var msg = map[string]string{
		"foo": "bar",
	}

	requestBody, _ := json.Marshal(msg)
	req := httptest.NewRequest(http.MethodPost, "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))

	var _, err = parseTelegramRequest(req)

	if err == nil {
		t.Error("Expected an error, got <nil>")
	}
}

func TestParseTelegramRequestErrorOnDecode(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer([]byte{1, 2}))

	var _, err = parseTelegramRequest(req)

	if err == nil {
		t.Error("Expected an error, got <nil>")
	}
}
