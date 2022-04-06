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

type sendMessagePayload struct {
	ChatID      int64               `json:"chat_id"`
	Text        string              `json:"text"`
	Entities    []sendMessageEntity `json:"entities,omitempty"`
	ReplyMarkup replyMarkup         `json:"reply_markup,omitempty"`
	ParseMode   string              `json:"parse_mode,omitempty"`
}

type replyMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard,omitempty"`
}

type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	CallbackData                 string `json:"callback_data"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
}

type sendMessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

type Client struct {
}

func (c *Client) SendMessage(botToken string, receiverID int64, message string, inlineButtons [][]InlineKeyboardButton) error {
	log.Println("SendMessage is called")
	payload := sendMessagePayload{
		ChatID:      receiverID,
		Text:        message,
		ReplyMarkup: replyMarkup{InlineKeyboard: inlineButtons},
	}
	data, err := json.Marshal(payload)

	log.Println("Request body: ", string(data))

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
