package models

import (
	"fmt"
	"log"
	"time"

	"example.com/chat/db"
)

type Message struct {
	ID int `json:"id"`
	ChannelID int `json:"channel_id"`
	UserID int `json:"user_id"`
	Content string `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	User User `json:"user"`
}

func SaveMessage(message Message) error{
	stmt, err := db.DB.Prepare("INSERT INTO messages (channel_id, user_id, content, timestamp) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(message.ChannelID, message.UserID, message.Content, message.Timestamp)
	return err
}

func GetMessagebyChannel(channelID int) ([]Message, error) {
	rows, err := db.DB.Query("SELECT id, channel_id, user_id, content, timestamp FROM messages WHERE channel_id = ?", channelID)
	if err != nil {
		log.Printf("Error querying messages: %v", err)
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ChannelID, &message.UserID, &message.Content, &message.Timestamp); err != nil {
			log.Printf("Error scanning message row: %v", err)
			return nil,  fmt.Errorf("error scanning message row: %w", err)
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return messages, nil
}

