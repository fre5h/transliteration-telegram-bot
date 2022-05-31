package handler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/fre5h/transliteration-telegram-bot/internal/model"

	"github.com/fre5h/transliteration-go"
)

func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	var result string

	var update, err = parseTelegramRequest(r)
	if nil != err {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	switch update.Message.Text {
	case "":
		result = "🤔 Вибачайте, але я вмію транслітерувати лише текстові повідомлення"
	case "/start":
		result = "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований текст 🇬🇧"
	default:
		result = transliteration.UkrToLat(update.Message.Text)
	}

	if responseBody, err := sendTextMessageToChat(update.Message.Chat.Id, result); err != nil {
		log.Printf("error %s from telegram, response body is %s", err.Error(), responseBody)
	}
}

func parseTelegramRequest(r *http.Request) (update *model.Update, err error) {
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())

		return nil, err
	}

	if 0 == update.UpdateId {
		log.Printf("invalid update id, got update id = 0")

		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}

	return update, nil
}

func sendTextMessageToChat(chatId int, text string) (string, error) {
	var botApiUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	response, errRequest := http.PostForm(
		botApiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if nil != errRequest {
		log.Printf("error when posting text to the chat: %s", errRequest.Error())

		return "", errRequest
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if nil != err {
			log.Printf("error on closing body: %s", err.Error())
		}
	}(response.Body)

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)

	if nil != errRead {
		log.Printf("error in parsing telegram answer %s", errRead.Error())

		return "", errRead
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}
