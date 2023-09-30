package telegram

import (
	"bot/clients/telegram"
	"bot/events"
	"bot/lib/errorH"
	"bot/storage"
	"errors"
)

var ErrUnknownEventType = errors.New("unknown event")
var ErrUnknownMetaAssertionType = errors.New("unknown meta type")

type Meta struct {
	ChatId   int
	Username string
}

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, errorH.Wrap("cant get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:

	}
}

func (p *Processor) processMessage(e events.Event) error {
	meta, err := meta(e)
	if err != nil {
		return errorH.Wrap("cant process message", err)
	}

}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, errorH.Wrap("cant assert meta type", ErrUnknownMetaAssertionType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	t := fetchType(upd)

	res := events.Event{
		Type: t,
		Text: fetchText(upd),
	}

	if t == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.Id,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
