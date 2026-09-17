package session

import (
	"errors"
	"sync"

	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/protocol"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

var (
	ErrAlreadyEstablished = errors.New("session is already established")
	ErrNotEstablished     = errors.New("session is not established")
)

type State string

const (
	StateNew         State = "new"
	StateEstablished State = "established"
	StateClosed      State = "closed"
)

type Session struct {
	mu       sync.Mutex
	state    State
	localID  identity.Identity
	remoteID string
	conn     *transport.Connection
}

func New(localID identity.Identity, conn *transport.Connection) (*Session, error) {
	if err := localID.Validate(); err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("connection cannot be nil")
	}

	return &Session{
		state:   StateNew,
		localID: localID,
		conn:    conn,
	}, nil
}

func (s *Session) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state
}

func (s *Session) RemoteID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.remoteID
}

func (s *Session) Establish(remoteID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == StateEstablished {
		return ErrAlreadyEstablished
	}

	if s.state == StateClosed {
		return ErrNotEstablished
	}

	if remoteID == "" {
		return errors.New("remote node ID cannot be empty")
	}

	if len(remoteID) != identity.NodeIDLength {
		return errors.New("invalid remote node ID")
	}

	s.remoteID = remoteID
	s.state = StateEstablished

	return nil
}

func (s *Session) Send(message protocol.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateEstablished {
		return ErrNotEstablished
	}

	return s.conn.Send(message)
}

func (s *Session) Receive() (protocol.Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != StateEstablished {
		return protocol.Message{}, ErrNotEstablished
	}

	return s.conn.Receive()
}

func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state == StateClosed {
		return nil
	}

	err := s.conn.Close()
	s.state = StateClosed

	return err
}
