package telegram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Client interface {
	SendTextMessageToChat(int, string) (string, error)
}

type HttpClient struct {
	baseUrl string
	token   string
}

func NewClient() *HttpClient {
	return &HttpClient{
		baseUrl: "https://api.telegram.org/bot",
		token:   os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (c HttpClient) SendTextMessageToChat(chatId int, text string) (string, error) {
	var botApiUrl = c.baseUrl + c.token + "/sendMessage"

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
