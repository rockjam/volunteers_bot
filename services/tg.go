package services

import (
	"bytes"
	"dv/adapters/httpbot"
	"dv/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	startCommand  = "/start"
	helpCommand   = "/help"
	commandPrefix = "/"
)

var (
	startCommandText = `
Greetings from volunteers_ua bot.

Available commands:
- /start, /help: show this message
- /info <Location>: show messages from volunteers about <Location>, i.e. /info Berlin
`
	invalidInfoCommandText = `
Invalid info command, maybe location is missing. Example: /info Berlin
`
	emptyInlineButtons [][]httpbot.InlineKeyboardButton
)

type Message struct {
	db       *gorm.DB
	botToken string
	botName  string
	chatID   int64
	httpBot  httpbot.Client
}

func NewMessage(db *gorm.DB, botToken string, botName string, chatID int64) Message {
	return Message{
		db:       db,
		chatID:   chatID,
		botToken: botToken,
		botName:  botName,
		httpBot:  httpbot.Client{},
	}
}

type cursor struct {
	location  string
	direction string
	timestamp int64
}

func newCursor(s string) (cursor, error) {
	parts := strings.Split(s, "§")

	if len(parts) != 3 {
		return cursor{}, errors.New("failed to parse the cursor")
	}

	location := parts[0]
	direction := parts[1]
	timestamp, err := strconv.ParseInt(parts[2], 0, 64)

	if err != nil {
		return cursor{}, err
	}

	return cursor{location: location, direction: direction, timestamp: timestamp}, nil
}

func (c *cursor) format() string {
	return fmt.Sprintf("%s§%s§%d", c.location, c.direction, c.timestamp)
}

func (m *Message) IncomingMessageHTTPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = m.HandleIncomingMessage(data)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (m *Message) HandleIncomingMessage(message []byte) error {

	log.Println("Incoming message: " + string(message))

	request := models.WebhookRequest{}
	err := json.Unmarshal(message, &request)
	if err != nil {
		return err
	}

	if request.Message != nil {
		log.Println("Message: ", request.Message)
		log.Println(fmt.Sprintf("message chatID is %d, default chat id is %d", request.Message.Chat.ID, m.chatID))
		if request.Message.Chat.ID == m.chatID {
			if strings.HasPrefix(request.Message.Text, fmt.Sprintf("%s@%s", startCommand, m.botName)) ||
				strings.HasPrefix(request.Message.Text, fmt.Sprintf("%s@%s", helpCommand, m.botName)) {
				return m.processStartCommand(request)
			} else if strings.HasPrefix(request.Message.Text, commandPrefix) {
				return m.processUnknownCommand(request)
			} else {
				return m.processGeneralMessage(request)
			}
		}

		if request.Message.Chat.Type == "private" {
			if strings.HasPrefix(request.Message.Text, startCommand) ||
				strings.HasPrefix(request.Message.Text, helpCommand) {
				return m.processStartCommand(request)
			} else if strings.HasPrefix(request.Message.Text, commandPrefix) {
				return m.processUnknownCommand(request)
			} else {
				return m.processInfoCommand(request)
			}
		}
	}

	if request.InlineQuery != nil {
		log.Println("Inline query: ", request.InlineQuery)
	}

	if request.CallbackQuery != nil {
		log.Println("Callback query: ", request.CallbackQuery)
		return m.processInfoCallback(request.CallbackQuery.Message.Chat.ID, request.CallbackQuery.Data)
	}

	return nil
}

func requestToMessages(req models.WebhookRequest) (models.Message, []models.MessageHashtag) {
	msg := models.Message{
		ID:            req.Message.MessageId,
		ChatID:        req.Message.Chat.ID,
		UpdateID:      req.UpdateID,
		FromID:        req.Message.From.ID,
		FromFirstName: req.Message.From.FirstName,
		FromLastName:  req.Message.From.LastName,
		FromUsername:  req.Message.From.Username,
		Timestamp:     req.Message.Date,
		Content:       req.Message.Text,
	}
	hashtags := make([]models.MessageHashtag, 0, len(req.Message.Entities))
	for _, entity := range req.Message.Entities {
		if entity.Type != "hashtag" {
			continue
		}
		left := entity.Offset
		right := entity.Offset + entity.Length
		if left > len(req.Message.Text) || right > len(req.Message.Text) {
			log.Println(fmt.Sprintf("invalid message entities: %v", req))
			continue
		}
		hashtags = append(hashtags, models.MessageHashtag{
			MessageID: msg.ID,
			ChatID:    msg.ChatID,
			Hashtag:   strings.ToLower(msg.Content[left:right]),
		})
	}
	return msg, hashtags
}

func (m *Message) storeMessage(msg models.Message, hashtags []models.MessageHashtag) error {
	if len(hashtags) == 0 {
		log.Println("no hash tags in message, returning")
		return nil
	}
	log.Println("storing message with id ", msg.ID)
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		for _, hashtag := range hashtags {
			if err := tx.Create(&hashtag).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *Message) processGeneralMessage(request models.WebhookRequest) error {
	return m.storeMessage(requestToMessages(request))
}

func (m *Message) processStartCommand(request models.WebhookRequest) error {
	return m.httpBot.SendMessage(m.botToken, request.Message.Chat.ID, startCommandText, emptyInlineButtons)
}

func (m *Message) processUnknownCommand(request models.WebhookRequest) error {
	message := "Unknown command. Send city/region/country to get a latest update, or /start for help"
	return m.httpBot.SendMessage(m.botToken, request.Message.Chat.ID, message, emptyInlineButtons)
}

func (m *Message) processInfoCallback(senderID int64, query string) error {
	cursor, err := newCursor(query)

	if err != nil {
		return err
	}

	return m.sendLocationResponse(senderID, cursor)
}

func (m *Message) processInfoCommand(request models.WebhookRequest) error {
	location := strings.TrimSpace(request.Message.Text)
	location = strings.Split(location, " ")[0]

	// TODO: find better way to put an initial timestamp, or find a way not to put timestamp and direction
	cursor := cursor{location: location, direction: "o", timestamp: math.MaxInt64}

	return m.sendLocationResponse(request.Message.Chat.ID, cursor)
}

func (m *Message) sendLocationResponse(senderID int64, c cursor) error {
	location := c.location
	if len(location) == 0 {
		return m.httpBot.SendMessage(m.botToken, senderID, invalidInfoCommandText, emptyInlineButtons)
	}
	results, err := m.getMessage(c)
	if err != nil {
		return err
	}
	log.Println("results: ", results)

	if results.message != nil {
		message := *results.message

		navButtons := make([]httpbot.InlineKeyboardButton, 2)

		if results.hasOlder {
			c := cursor{location: location, direction: "o", timestamp: message.Timestamp}
			b := httpbot.InlineKeyboardButton{
				Text:         "Older",
				CallbackData: c.format(),
			}
			navButtons = append(navButtons, b)
		}
		if results.hasNewer {
			c := cursor{location: location, direction: "n", timestamp: message.Timestamp}
			b := httpbot.InlineKeyboardButton{
				Text:         "Newer",
				CallbackData: c.format(),
			}
			navButtons = append(navButtons, b)
		}
		inlineButtons := [][]httpbot.InlineKeyboardButton{navButtons}
		return m.httpBot.SendMessage(m.botToken, senderID, formatMessage(message), inlineButtons)
	} else {
		return m.httpBot.SendMessage(m.botToken, senderID, fmt.Sprintf("No messages for location '%s' were found", location), emptyInlineButtons)
	}
}

func formatMessage(message models.Message) string {
	outputBuf := new(bytes.Buffer)
	messageTS := time.Unix(message.Timestamp, 0)
	outputBuf.WriteString(fmt.Sprintf("%s from @%s (%s %s)\n",
		messageTS.Format(time.RFC822), message.FromUsername, message.FromFirstName, message.FromLastName))
	outputBuf.WriteString(message.Content)
	return outputBuf.String()
}

type searchResult struct {
	message  *models.Message
	hasOlder bool
	hasNewer bool
}

func (m *Message) getMessage(c cursor) (searchResult, error) {
	searchRes := searchResult{}

	tag := "#" + strings.ToLower(c.location)

	message, err := m.fetchMessages(tag, c)
	if err != nil {
		return searchRes, err
	}

	if message != nil {
		hasNewer, err := m.hasNewerMessages(*message, tag)
		if err != nil {
			return searchRes, err
		}

		hasOlder, err := m.hasOlderMessages(*message, tag)
		if err != nil {
			return searchRes, err
		}

		searchRes.message = message
		searchRes.hasNewer = hasNewer
		searchRes.hasOlder = hasOlder
	}

	return searchRes, nil
}

func (m *Message) fetchMessages(tag string, c cursor) (*models.Message, error) {
	timestamp := c.timestamp
	var direction string
	var sorting string
	if c.direction == "o" {
		direction = "<"
		sorting = "desc"
	} else {
		direction = ">"
		sorting = "asc"
	}
	log.Println("getting messages for tag: ", tag, "timestamp: ", timestamp, "direction: ", direction)

	var messages []models.Message

	messageQuery := fmt.Sprintf(`
SELECT m.*
FROM messages m
         JOIN message_hashtags h ON m.id = h.message_id AND m.chat_id = h.chat_id AND h.hashtag = ?
where m.timestamp %s ?
ORDER BY m.timestamp %s
LIMIT 1;`, direction, sorting)
	res := m.db.Raw(messageQuery, tag, timestamp).Scan(&messages)

	if len(messages) == 1 {
		return &messages[0], res.Error
	} else {
		return nil, res.Error
	}
}

func (m *Message) hasNewerMessages(message models.Message, tag string) (bool, error) {
	var hasNewer bool

	res := m.db.Raw(`
SELECT count(*) > 0
FROM messages m
         JOIN message_hashtags h ON m.id = h.message_id AND m.chat_id = h.chat_id AND h.hashtag = ?
where m.timestamp > ?`, tag, message.Timestamp).Scan(&hasNewer)

	return hasNewer, res.Error
}

func (m *Message) hasOlderMessages(message models.Message, tag string) (bool, error) {
	var hasOlder bool

	res := m.db.Raw(`
SELECT count(*) > 0
FROM messages m
         JOIN message_hashtags h ON m.id = h.message_id AND m.chat_id = h.chat_id AND h.hashtag = ?
where m.timestamp < ?`, tag, message.Timestamp).Scan(&hasOlder)

	return hasOlder, res.Error
}
