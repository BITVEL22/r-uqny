package transport

import (
	"net"
	"testing"

	"github.com/BITVEL22/r-uqny/internal/protocol"
)

func TestListenTCP(t *testing.T) {
	listener, err := ListenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	defer listener.Close()

	if listener.Address() == nil {
		t.Fatal("listener address should not be nil")
	}

	done := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()

		done <- nil
	}()

	conn, err := net.Dial("tcp", listener.Address().String())
	if err != nil {
		t.Fatalf("net.Dial() error = %v", err)
	}
	defer conn.Close()

	if err := <-done; err != nil {
		t.Fatalf("Accept() error = %v", err)
	}
}

func TestDialTCP(t *testing.T) {
	listener, err := ListenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	defer listener.Close()

	done := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()

		done <- nil
	}()

	conn, err := DialTCP(listener.Address().String())
	if err != nil {
		t.Fatalf("DialTCP() error = %v", err)
	}
	defer conn.Close()

	if err := <-done; err != nil {
		t.Fatalf("Accept() error = %v", err)
	}
}

func TestSendReceive(t *testing.T) {
	listener, err := ListenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	defer listener.Close()

	done := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()

		buffer := make([]byte, 1024)
		n, err := Receive(conn, buffer)
		if err != nil {
			done <- err
			return
		}

		if string(buffer[:n]) != "hello r/uqny" {
			done <- &unexpectedDataError{}
			return
		}

		done <- nil
	}()

	conn, err := DialTCP(listener.Address().String())
	if err != nil {
		t.Fatalf("DialTCP() error = %v", err)
	}
	defer conn.Close()

	if err := Send(conn, []byte("hello r/uqny")); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("server error = %v", err)
	}
}

type unexpectedDataError struct{}

func (*unexpectedDataError) Error() string {
	return "received unexpected data"
}

func TestSendReceiveProtocolMessage(t *testing.T) {
	listener, err := ListenTCP("127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	defer listener.Close()

	message := protocol.NewTextMessage(
		"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"test-message-1",
		"hello from r/uqny",
	)

	data, err := protocol.Encode(message)
	if err != nil {
		t.Fatalf("protocol.Encode() error = %v", err)
	}

	done := make(chan error, 1)

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			done <- err
			return
		}
		defer conn.Close()

		receivedData, err := ReceiveMessage(conn)
		if err != nil {
			done <- err
			return
		}

		receivedMessage, err := protocol.Decode(receivedData)
		if err != nil {
			done <- err
			return
		}

		if receivedMessage.MessageID != message.MessageID {
			done <- &unexpectedDataError{}
			return
		}

		if receivedMessage.Payload != message.Payload {
			done <- &unexpectedDataError{}
			return
		}

		done <- nil
	}()

	conn, err := DialTCP(listener.Address().String())
	if err != nil {
		t.Fatalf("DialTCP() error = %v", err)
	}
	defer conn.Close()

	if err := SendMessage(conn, data); err != nil {
		t.Fatalf("SendMessage() error = %v", err)
	}

	if err := <-done; err != nil {
		t.Fatalf("server error = %v", err)
	}
}
