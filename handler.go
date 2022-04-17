package transliteration_telegram_bot_handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/fre5h/transliteration-go"
)

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	var result string
	if "" == update.Message.Text {
		result = "🤔 Вибачайте, але я вмію траслітерувати лише текстові повідомлення"
	} else if "/start" == update.Message.Text {
		result = "Просто напишіть мені текст на українській мові 🇺🇦 і у відповідь отримаєте транслітерований текст 🇬🇧"
	} else {
		result = transliteration.UkrToLat(update.Message.Text)
	}

	if responseBody, err := sendTextMessageToChat(update.Message.Chat.Id, result); err != nil {
		log.Printf("error %s from telegram, reponse body is %s", err.Error(), responseBody)
	}
}

func sendTextMessageToChat(chatId int, text string) (string, error) {
	var borApiUrl = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	response, errRequest := http.PostForm(
		borApiUrl,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})
	defer response.Body.Close()

	if errRequest != nil {
		log.Printf("error when posting text to the chat: %s", errRequest.Error())

		return "", errRequest
	}

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())

		return "", errRead
	}

	bodyString := string(bodyBytes)

	return bodyString, nil
}

func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())

		return nil, err
	}

	if update.UpdateId == 0 {
		log.Printf("invalid update id, got update id = 0")

		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}

	return &update, nil
}
