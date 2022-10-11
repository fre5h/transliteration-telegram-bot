package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/fre5h/transliteration-go"

	"github.com/fre5h/transliteration-telegram-bot/internal/model"
)

func HandleTelegramWebHook(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	var update model.Update

	err := json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return events.LambdaFunctionURLResponse{
				StatusCode: 500,
				Body:       err.Error(),
			},
			err
	}

	if 0 == update.UpdateId {
		errorMessage := "invalid update id of 0 indicates failure to parse incoming update"

		return events.LambdaFunctionURLResponse{
				StatusCode: 400,
				Body:       errorMessage,
			},
			errors.New(errorMessage)
	}

	var result string

	switch update.Message.Text {
	case "":
		result = "🤔 Вибачайте, але я вмію транслітерувати лише текстові повідомлення"
	case "/start":
		result = "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований 🇬🇧 текст"
	default:
		result = transliteration.UkrToLat(update.Message.Text)
	}

	if responseBody, err := sendTextMessageToChat(update.Message.Chat.Id, result); err != nil {
		return events.LambdaFunctionURLResponse{
				Body:       err.Error(),
				StatusCode: 500,
			},
			fmt.Errorf("error %s from telegram, response body is %s", err.Error(), responseBody)
	}

	return events.LambdaFunctionURLResponse{StatusCode: 200, Body: "OK"}, nil
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

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)

	if nil != errRead {
		return "", fmt.Errorf("error in parsing telegram answer %s", errRead.Error())
	}

	return string(bodyBytes), nil
}
