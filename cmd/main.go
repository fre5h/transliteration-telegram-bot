package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/fre5h/transliteration-telegram-bot/internal/handler"
	"github.com/fre5h/transliteration-telegram-bot/internal/telegram"
)

func main() {
	telegramClient := telegram.NewClient()
	lambdaHandler := handler.NewLambdaHandler(*telegramClient)
	lambda.Start(lambdaHandler.HandleLambdaRequest)
}
