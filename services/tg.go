package services

import (
	"bytes"
	"dv/adapters/httpbot"
	"dv/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	defaultChatID  = int64(-1001722078663)
	maxGetMessages = 10
)

type Message struct {
	db       *gorm.DB
	botToken string
	chatID   int64
	httpBot  httpbot.Client
}

func NewMessage(db *gorm.DB, botToken string) Message {
	chatID, err := strconv.ParseInt(os.Getenv("GROUP_CHAT_ID"), 10, 64)
	if err != nil {
		chatID = defaultChatID
	}
	return Message{
		db:       db,
		chatID:   chatID,
		botToken: botToken,
		httpBot:  httpbot.Client{},
	}
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

	log.Println(fmt.Sprintf("message chatID is %d, default chat id is %d", request.Message.Chat.ID, m.chatID))
	if request.Message.Chat.ID == m.chatID {
		err = m.processGeneralMessage(request)
		if err != nil {
			return err
		}
		return nil
	}
	if request.Message.Chat.Type == "private" {
		err = m.processBotCommand(request)
		if err != nil {
			return err
		}
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

//query := db.Table("order").Select("MAX(order.finished_at) as latest").Joins("left join user user on order.user_id = user.id").Where("user.age > ?", 18).Group("order.user_id")
//db.Model(&Order{}).Joins("join (?) q on order.finished_at = q.latest", query).Scan(&results)

func (m *Message) processGeneralMessage(request models.WebhookRequest) error {
	return m.storeMessage(requestToMessages(request))
}

const (
	startCommand = "/start"
	helpCommand  = "/help"
	infoCommand  = "/info"
)

func (m *Message) processBotCommand(request models.WebhookRequest) error {
	log.Println("processing bot command ", request.Message.Text)
	switch {
	case strings.HasPrefix(request.Message.Text, startCommand):
		return m.processStartCommand(request.Message.From.ID)
	case strings.HasPrefix(request.Message.Text, helpCommand):
		return m.processStartCommand(request.Message.From.ID)
	case strings.HasPrefix(request.Message.Text, infoCommand):
		return m.processInfoCommand(request.Message.From.ID, request.Message.Text)
	}
	return nil
}

func (m *Message) processStartCommand(senderID int64) error {
	log.Println("processing start command")
	return m.httpBot.SendWelcomeMessage(m.botToken, senderID)
}

func (m *Message) processInfoCommand(senderID int64, command string) error {
	location := strings.TrimPrefix(command, infoCommand)
	location = strings.TrimPrefix(location, " ")
	location = strings.Split(location, " ")[0]
	if len(location) == 0 {
		return m.httpBot.SendInvalidInfoCommandMessage(m.botToken, senderID)
	}
	messages, err := m.getMessagesByTag(location)
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return m.httpBot.SendCustomMessage(m.botToken, senderID,
			fmt.Sprintf("No messages for location '%s' were found", location))
	}
	outputBuf := new(bytes.Buffer)
	outputBuf.WriteString("Following messages were found:\n")
	for _, message := range messages {
		outputBuf.WriteString("-----------------------------------\n")
		messageTS := time.Unix(message.Timestamp, 0)
		outputBuf.WriteString(fmt.Sprintf("From @%s (%s %s) at %s\n",
			message.FromUsername, message.FromFirstName, message.FromLastName, messageTS.Format(time.RFC3339)))
		outputBuf.WriteString(message.Content + "\n")
	}
	return m.httpBot.SendCustomMessage(m.botToken, senderID, outputBuf.String())
}

func (m *Message) getMessagesByTag(tag string) ([]models.Message, error) {
	tag = "#" + strings.ToLower(tag)
	log.Println("getting messages for tag ", tag)
	var results = make([]models.Message, 0, maxGetMessages)

	res := m.db.Raw(`
SELECT m.* FROM messages m 
JOIN message_hashtags h ON m.id = h.message_id AND m.chat_id = h.chat_id AND h.hashtag = ? 
ORDER BY m.timestamp desc LIMIT ?;`, tag, maxGetMessages).Scan(&results)
	if res.Error != nil {
		return nil, res.Error
	}
	return results, nil
}
