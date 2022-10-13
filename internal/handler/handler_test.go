package handler

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

// func TestParseTelegramRequest(t *testing.T) {
// 	var update = model.Update{
// 		UpdateId: 1,
// 		Message: model.Message{
// 			Text: "привіт",
// 			Chat: model.Chat{
// 				Id: 1,
// 			},
// 		},
// 	}
//
// 	requestBody, err := json.Marshal(update)
// 	if err != nil {
// 		t.Errorf("Failed to marshal update in json, got %s", err.Error())
// 	}
// 	req := httptest.NewRequest(http.MethodPost, "https://myTelegramWebHookHandler.com/secretToken", bytes.NewBuffer(requestBody))
//
// 	var updateToTest, errParse = parseTelegramRequest(req)
// 	if errParse != nil {
// 		t.Errorf("Expected a <nil> error, got %s", errParse.Error())
// 	}
//
// 	if *updateToTest != update {
// 		t.Errorf("Expected update %v, got %v", update, updateToTest)
// 	}
// }

type MockHttpClient struct {
	baseUrl string
	token   string
}

func NewMockHttpClient() *MockHttpClient {
	return &MockHttpClient{
		baseUrl: "https://api.telegram.org/bot",
		token:   os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (c MockHttpClient) SendTextMessageToChat(_ int, _ string) (string, error) {
	return "OK", nil
}

func TestParseTelegramRequestInvalid(t *testing.T) {
	handler := NewLambdaHandler(NewMockHttpClient())
	request := events.LambdaFunctionURLRequest{Body: "{\"update_id\":0}"}
	var response, _ = handler.HandleLambdaRequest(request)

	if response.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", response.StatusCode)
	}
}

func TestParseTelegramRequestErrorOnDecode(t *testing.T) {
	handler := NewLambdaHandler(NewMockHttpClient())
	request := events.LambdaFunctionURLRequest{Body: "123"}
	var response, _ = handler.HandleLambdaRequest(request)

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
