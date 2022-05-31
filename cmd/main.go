package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fre5h/transliteration-telegram-bot/internal/handler"
)

func main() {
	http.HandleFunc("/handler/"+os.Getenv("WEB_SOCKET_SECRET"), handler.HandleTelegramWebHook)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
