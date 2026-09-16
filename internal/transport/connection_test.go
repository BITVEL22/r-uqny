package transport

import (
	"testing"

	"github.com/BITVEL22/r-uqny/internal/protocol"
)

func TestConnectionSendReceive(t *testing.T) {
	listener, err := ListenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	defer listener.Close()

	message := protocol.NewTextMessage(
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"connection-test-1",
		"hello through Connection",
	)

	done := make(chan error, 1)

	go func() {
		rawConn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}

		conn := NewConnection(rawConn)
		defer conn.Close()

		received, err := conn.Receive()
		if err != nil {
			done <- err
			return
		}

		if received.MessageID != message.MessageID {
			done <- &unexpectedDataError{}
			return
		}

		if received.Payload != message.Payload {
			done <- &unexpectedDataError{}
			return
		}

		done <- nil
	}()

	rawConn, err := DialTCP(listener.Address().String())
	if err != nil {
		t.Fatalf("DialTCP() error = %v", err)
	}

	conn := NewConnection(rawConn)
	defer conn.Close()

	if err := conn.Send(message); err != nil {
		t.Fatalf("Connection.Send() error = %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("server error = %v", err)
	}
}
