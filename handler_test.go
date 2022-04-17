package transliteration_telegram_bot_handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestParseTelegramRequest(t *testing.T) {
	var update = Update{
		UpdateId: 1,
		Message: Message{
			Text: "привіт",
			Chat: Chat{
				Id: 1,
			},
		},
	}

	requestBody, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Failed to marshal update in json, got %s", err.Error())
	}
	req := httptest.NewRequest("POST", "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))

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
	req := httptest.NewRequest("POST", "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))

	var _, err = parseTelegramRequest(req)

	if err == nil {
		t.Error("Expected an error, got <nil>")
	}
}
