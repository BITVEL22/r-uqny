package transport

import (
	"encoding/binary"
	"io"
	"net"
)

type TCPListener struct {
	listener net.Listener
}

func ListenTCP(address string) (*TCPListener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return &TCPListener{
		listener: listener,
	}, nil
}

func (t *TCPListener) Accept() (net.Conn, error) {
	return t.listener.Accept()
}

func (t *TCPListener) Close() error {
	return t.listener.Close()
}

func (t *TCPListener) Address() net.Addr {
	return t.listener.Addr()
}

func DialTCP(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

func Send(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	return err
}

func Receive(conn net.Conn, buffer []byte) (int, error) {
	return conn.Read(buffer)
}

func SendMessage(conn net.Conn, data []byte) error {
	length := uint32(len(data))

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)

	if _, err := conn.Write(header); err != nil {
		return err
	}

	_, err := conn.Write(data)
	return err
}

func ReceiveMessage(conn net.Conn) ([]byte, error) {
	header := make([]byte, 4)

	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)

	data := make([]byte, length)

	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return data, nil
}
