package telegram

import (
	"net/http"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type UpdatesResponse struct {
	Description string   `json:"description"`
	Ok          bool     `json:"ok"`
	Result      []Update `json:"result"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

type From struct {
	Username string `json:"username"`
}

type Chat struct {
	Id int `json:"id"`
}
