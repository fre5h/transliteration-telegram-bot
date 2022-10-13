package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/fre5h/transliteration-go"

	"github.com/fre5h/transliteration-telegram-bot/internal/model"
)

func HandleTelegramWebHook(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var update model.Update

	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return createLambdaResponse(http.StatusInternalServerError, "Error on unmarshaling json")
	}

	if 0 == update.UpdateId {
		return createLambdaResponse(http.StatusBadRequest, "Update id of 0 indicates failure to parse incoming update")
	}

	if responseBody, err := sendTextMessageToChat(update.Message.Chat.Id, prepareResult(update.Message.Text)); err != nil {
		log.Printf("error %s from telegram, response body is %s", err.Error(), responseBody)

		return createLambdaResponse(http.StatusInternalServerError, "Error on request to Telegram")
	}

	return createLambdaResponse(http.StatusOK, "OK")
}

func createLambdaResponse(statusCode int, body string) (events.LambdaFunctionURLResponse, error) {
	return events.LambdaFunctionURLResponse{StatusCode: statusCode, Body: body}, nil
}

func prepareResult(text string) (result string) {
	switch text {
	case "":
		result = "🤔 Вибачайте, але я вмію транслітерувати лише текстові повідомлення"
	case "/start":
		result = "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований 🇬🇧 текст"
	default:
		result = transliteration.UkrToLat(text)
	}

	return result
}

func sendTextMessageToChat(chatId int, text string) (string, error) {
	var botApiUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	response, errRequest := http.PostForm(
		botApiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		},
	)

	if nil != errRequest {
		return "", fmt.Errorf("error when posting text to the chat: %s", errRequest.Error())
	}

	if 200 != response.StatusCode {
		return "", fmt.Errorf("status code of response is: %d", response.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if nil != err {
			log.Printf("error on closing body: %s", err.Error())
		}
	}(response.Body)

	var bodyBytes, errRead = io.ReadAll(response.Body)

	if nil != errRead {
		return "", fmt.Errorf("error in parsing telegram answer %s", errRead.Error())
	}

	return string(bodyBytes), nil
}
