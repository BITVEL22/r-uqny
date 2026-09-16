package transport

import (
	"net"

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

func (c *Connection) Close() error {
	return c.conn.Close()
}
