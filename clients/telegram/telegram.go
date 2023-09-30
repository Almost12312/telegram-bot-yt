package telegram

import (
	"bot/lib/errorH"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdates  = "getUpdates"
	sendMessage = "sendMessage"
	cantSendReq = "Can`t do request"
)

func New(host string, token string) Client {
	return Client{
		client:   http.Client{},
		basePath: basePath(token),
		host:     host,
	}
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chatID", strconv.Itoa(chatId))
	q.Add("chatID", text)

	_, err := c.doRequest(sendMessage, q)
	if err != nil {
		return errorH.Wrap(cantSendReq, err)
	}

	return nil
}

func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdates, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, err
}

func (c *Client) doRequest(endPoint string, query url.Values) (data []byte, err error) {

	defer func() {
		err = errorH.WrapIfErr(cantSendReq, err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, endPoint),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func basePath(token string) string {
	return "bot" + token
}
