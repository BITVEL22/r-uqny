package protocol

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
)

const CurrentVersion = 1

const (
	TypeText = "text"
)

type Message struct {
	Version   int       `json:"version"`
	MessageID string    `json:"message_id"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}

func NewTextMessage(sender, messageID, payload string) Message {
	return Message{
		Version:   CurrentVersion,
		MessageID: messageID,
		Type:      TypeText,
		Sender:    sender,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

func (m Message) Validate() error {
	if m.Version != CurrentVersion {
		return errors.New("unsupported protocol version")
	}

	if m.MessageID == "" {
		return errors.New("message ID cannot be empty")
	}

	if m.Type == "" {
		return errors.New("message type cannot be empty")
	}

	if m.Sender == "" {
		return errors.New("sender cannot be empty")
	}

	if len(m.Sender) != identity.NodeIDLength {
		return errors.New("invalid sender node ID")
	}

	if m.Timestamp.IsZero() {
		return errors.New("timestamp cannot be zero")
	}

	return nil
}

func Encode(m Message) ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}

	return json.Marshal(m)
}

func Decode(data []byte) (Message, error) {
	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		return Message{}, err
	}

	if err := message.Validate(); err != nil {
		return Message{}, err
	}

	return message, nil
}
