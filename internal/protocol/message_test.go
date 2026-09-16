package protocol

import (
	"testing"
	"time"
)

const testSender = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func TestNewTextMessage(t *testing.T) {
	message := NewTextMessage(
		testSender,
		"msg-001",
		"Hello r/uqny!",
	)

	if message.Version != CurrentVersion {
		t.Fatalf("expected version %d, got %d", CurrentVersion, message.Version)
	}

	if message.MessageID != "msg-001" {
		t.Fatalf("unexpected message ID: %s", message.MessageID)
	}

	if message.Type != TypeText {
		t.Fatalf("expected type %q, got %q", TypeText, message.Type)
	}

	if message.Sender != testSender {
		t.Fatalf("unexpected sender: %s", message.Sender)
	}

	if message.Payload != "Hello r/uqny!" {
		t.Fatalf("unexpected payload: %s", message.Payload)
	}

	if message.Timestamp.IsZero() {
		t.Fatal("expected timestamp to be set")
	}
}

func TestMessageValidate(t *testing.T) {
	message := NewTextMessage(testSender, "msg-001", "hello")

	if err := message.Validate(); err != nil {
		t.Fatalf("valid message should pass validation: %v", err)
	}
}

func TestInvalidMessage(t *testing.T) {
	message := NewTextMessage(testSender, "msg-001", "hello")

	tests := []struct {
		name   string
		mutate func(*Message)
	}{
		{
			name: "invalid version",
			mutate: func(m *Message) {
				m.Version = 999
			},
		},
		{
			name: "empty message ID",
			mutate: func(m *Message) {
				m.MessageID = ""
			},
		},
		{
			name: "empty type",
			mutate: func(m *Message) {
				m.Type = ""
			},
		},
		{
			name: "empty sender",
			mutate: func(m *Message) {
				m.Sender = ""
			},
		},
		{
			name: "invalid sender",
			mutate: func(m *Message) {
				m.Sender = "invalid"
			},
		},
		{
			name: "zero timestamp",
			mutate: func(m *Message) {
				m.Timestamp = time.Time{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			invalid := message
			test.mutate(&invalid)

			if err := invalid.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	original := NewTextMessage(
		testSender,
		"msg-001",
		"Hello r/uqny!",
	)

	data, err := Encode(original)
	if err != nil {
		t.Fatalf("failed to encode message: %v", err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("failed to decode message: %v", err)
	}

	if decoded.Version != original.Version {
		t.Fatal("version changed after encode/decode")
	}

	if decoded.MessageID != original.MessageID {
		t.Fatal("message ID changed after encode/decode")
	}

	if decoded.Type != original.Type {
		t.Fatal("message type changed after encode/decode")
	}

	if decoded.Sender != original.Sender {
		t.Fatal("sender changed after encode/decode")
	}

	if !decoded.Timestamp.Equal(original.Timestamp) {
		t.Fatal("timestamp changed after encode/decode")
	}

	if decoded.Payload != original.Payload {
		t.Fatal("payload changed after encode/decode")
	}
}

func TestDecodeInvalidJSON(t *testing.T) {
	if _, err := Decode([]byte("{invalid")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
