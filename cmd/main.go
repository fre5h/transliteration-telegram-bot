package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fre5h/transliteration-telegram-bot/internal/handler"
)

func main() {
	lambda.Start(handler.HandleTelegramWebHook)
}
