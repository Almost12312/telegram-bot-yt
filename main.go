package main

import (
	"flag"
	"log"

	"bot/clients/telegram"
)

const (
	host = "api.telegram.org"
)

func main() {

	tgClient := telegram.New(host, mustToken())

	_ = tgClient
}

func mustToken() string {
	token := flag.String("bot-token",
		"",
		"for start bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("No token for bot!")
	}

	return *token
}
