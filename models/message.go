package models

type WebhookRequest struct {
	UpdateID int64          `json:"update_id"`
	Message  WebhookMessage `json:"message"`
}

type WebhookMessage struct {
	MessageId int64 `json:"message_id"`
	From      struct {
		ID        int64  `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	} `json:"from"`
	Chat struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
	} `json:"chat"`
	Date     int64  `json:"date"`
	Text     string `json:"text"`
	Entities []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
}

type Message struct {
	ID            int64 `gorm:"primaryKey"`
	ChatID        int64 `gorm:"primaryKey"`
	UpdateID      int64
	FromID        int64
	FromUsername  string
	FromFirstName string
	FromLastName  string
	Timestamp     int64
	Content       string
}

type MessageHashtag struct {
	MessageID int64  `gorm:"primaryKey"`
	ChatID    int64  `gorm:"primaryKey"`
	Hashtag   string `gorm:"primaryKey"`
}
