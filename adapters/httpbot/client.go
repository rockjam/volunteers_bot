package httpbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	startCommandText = `
Greetings from volunteers_ua bot.

Available commands:
- /start, /help: show this message
- /info <Location>: show messages from volunteers about <Location>, i.e. /info Berlin
`
	startCommandEntities = []sendMessageEntity{
		{
			Type:   "bold",
			Offset: 16,
			Length: 13,
		},
		{
			Type:   "code",
			Offset: 92,
			Length: 16,
		},
		{
			Type:   "bot_command",
			Offset: 192,
			Length: 12,
		},
	}

	invalidInfoCommandText = `
Invalid info command, maybe location is missing. Example: /info Berlin
`
)

type sendMessagePayload struct {
	ChatID    int64               `json:"chat_id"`
	Text      string              `json:"text"`
	Entities  []sendMessageEntity `json:"entities,omitempty"`
	ParseMode string              `json:"parse_mode,omitempty"`
}

type sendMessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Client struct {
}

func (c *Client) SendWelcomeMessage(botToken string, receiverID int64) error {
	log.Println("SendWelcomeMessage is called")
	payload := sendMessagePayload{
		ChatID:   receiverID,
		Text:     startCommandText,
		Entities: startCommandEntities,
	}
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}
	u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	return c.sendMessage(data, u)
}

func (c *Client) SendInvalidInfoCommandMessage(botToken string, receiverID int64) error {
	log.Println("SendInvalidInfoCommandMessage is called")
	payload := sendMessagePayload{
		ChatID:   receiverID,
		Text:     invalidInfoCommandText,
		Entities: nil,
	}
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}
	u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	return c.sendMessage(data, u)
}

func (c *Client) SendCustomMessage(botToken string, receiverID int64, message string) error {
	log.Println("SendCustomMessage is called")
	payload := sendMessagePayload{
		ChatID: receiverID,
		Text:   message,
	}
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}
	u := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	return c.sendMessage(data, u)
}

func (c *Client) sendMessage(data []byte, u string) error {

	req, err := http.NewRequest(
		http.MethodPost,
		u,
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		log.Println("resp body: " + string(body))
		log.Println("resp status code: " + strconv.Itoa(resp.StatusCode))
		return errors.New("Invalid request status code " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
