package handler

import (
	"testing"

	"github.com/fre5h/transliteration-telegram-bot/internal/mocks"
)

func TestSendTextMessageToChatSuccessfully(t *testing.T) {
	client := NewTelegramHttpClient(mocks.NewMockHttpClient())
	response, err := client.SendTextMessageToChat(1, "test")

	if err != nil {
		t.Errorf("Expected no error, but \"%s\" give", err)
	}

	if response != "OK" {
		t.Errorf("Expected response to be \"OK\", got %s", response)
	}
}
