package main

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

	"github.com/fre5h/transliteration-go"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

func main() {
	http.HandleFunc("/handler/"+os.Getenv("WEB_SOCKET_SECRET"), HandleTelegramWebHook)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func HandleTelegramWebHook(_ http.ResponseWriter, r *http.Request) {
	var result string

	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

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

	if errRequest != nil {
		log.Printf("error when posting text to the chat: %s", errRequest.Error())

		return "", errRequest
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error on closing body: %s", err.Error())
		}
	}(response.Body)

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
