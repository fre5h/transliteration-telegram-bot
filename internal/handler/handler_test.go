package handler

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandleLambdaRequestInvalidUpdateId(t *testing.T) {
	handler := NewLambdaHandler(NewMockOkHttpClient())
	response, _ := handler.HandleLambdaRequest(events.LambdaFunctionURLRequest{Body: "{\"update_id\":0}"})

	if response.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", response.StatusCode)
	}
}

func TestHandleLambdaRequestErrorOnUnmarshal(t *testing.T) {
	handler := NewLambdaHandler(NewMockOkHttpClient())
	response, _ := handler.HandleLambdaRequest(events.LambdaFunctionURLRequest{Body: "123"})

	if response.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", response.StatusCode)
	}
}

func TestHandleLambdaRequestSuccessfully(t *testing.T) {
	handler := NewLambdaHandler(NewMockOkHttpClient())
	response, _ := handler.HandleLambdaRequest(events.LambdaFunctionURLRequest{Body: `{"update_id":1,"message":{"text":"привіт","chat":{"id":1}}}`})

	if response.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", response.StatusCode)
	}
}

func TestHandleLambdaRequestUnsuccessfullyRequestToTelegram(t *testing.T) {
	handler := NewLambdaHandler(NewMockFailedHttpClient())
	response, _ := handler.HandleLambdaRequest(events.LambdaFunctionURLRequest{Body: `{"update_id":1,"message":{"text":"привіт","chat":{"id":1}}}`})

	if response.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", response.StatusCode)
	}
}

func TestPrepareResult(t *testing.T) {
	if prepareResult("") != "🤔 Вибачайте, але я вмію транслітерувати лише текстові повідомлення" {
		t.Error("Expected another result for empty string")
	}

	if prepareResult("/start") != "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований 🇬🇧 текст" {
		t.Error("Expected another result for /start command")
	}

	result := prepareResult("тест")
	if result != "test" {
		t.Errorf("Expected \"test\", got %s", result)
	}
}
