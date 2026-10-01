package transport

import (
	"net"
	"time"

	"github.com/BITVEL22/r-uqny/internal/protocol"
)

type Connection struct {
	conn net.Conn
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) Send(message protocol.Message) error {
	data, err := protocol.Encode(message)
	if err != nil {
		return err
	}

	return SendMessage(c.conn, data)
}

func (c *Connection) Receive() (protocol.Message, error) {
	data, err := ReceiveMessage(c.conn)
	if err != nil {
		return protocol.Message{}, err
	}

	return protocol.Decode(data)
}

func (c *Connection) SendHandshake(handshake protocol.Handshake) error {
	data, err := handshake.Encode()
	if err != nil {
		return err
	}

	return SendMessage(c.conn, data)
}

func (c *Connection) ReceiveHandshake() (protocol.Handshake, error) {
	data, err := ReceiveMessage(c.conn)
	if err != nil {
		return protocol.Handshake{}, err
	}

	return protocol.DecodeHandshake(data)
}

func (c *Connection) Close() error {
	return c.conn.Close()
}

// SetDeadline sets the read and write deadlines on the connection.
func (c *Connection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the read deadline on the connection.
func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline on the connection.
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// RemoteAddr returns the remote address of the connection.
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// LocalAddr returns the local address of the connection.
func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
