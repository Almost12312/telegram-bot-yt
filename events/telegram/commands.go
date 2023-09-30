package telegram

import (
	"bot/lib/errorH"
	"bot/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/random"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) {
	text := strings.TrimSpace(text)

	log.Printf("command %s from user %s\n", text, username)

	if isAddCmd(text) {
	}

	switch text {
	case StartCmd:

	case HelpCmd:

	case RndCmd:

	default:

	}
}

func (p *Processor) savePage(pageUrl string, chatId int, username string) (err error) {
	defer func() { err = errorH.WrapIfErr("cant save page", err) }()

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatId)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
