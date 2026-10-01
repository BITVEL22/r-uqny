package session

import (
	"errors"
	"sync"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/protocol"
	"github.com/BITVEL22/r-uqny/internal/transport"
)

const HandshakeTimeout = 30 * time.Second

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
	stateMu  sync.RWMutex
	writeMu  sync.Mutex
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
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	return s.state
}

func (s *Session) RemoteID() string {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()

	return s.remoteID
}

func (s *Session) Establish(remoteID string) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

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

func (s *Session) HandshakeAsClient() error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state == StateEstablished {
		return ErrAlreadyEstablished
	}

	if s.state == StateClosed {
		return ErrNotEstablished
	}

	s.conn.SetDeadline(time.Now().Add(HandshakeTimeout))
	defer s.conn.SetDeadline(time.Time{})

	handshake, err := protocol.NewHandshake(s.localID)
	if err != nil {
		return err
	}

	if err := s.conn.SendHandshake(handshake); err != nil {
		return err
	}

	remoteHandshake, err := s.conn.ReceiveHandshake()
	if err != nil {
		return err
	}

	if err := remoteHandshake.Validate(); err != nil {
		return err
	}

	s.remoteID = remoteHandshake.NodeID
	s.state = StateEstablished

	return nil
}

func (s *Session) HandshakeAsServer() error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state == StateEstablished {
		return ErrAlreadyEstablished
	}

	if s.state == StateClosed {
		return ErrNotEstablished
	}

	s.conn.SetDeadline(time.Now().Add(HandshakeTimeout))
	defer s.conn.SetDeadline(time.Time{})

	remoteHandshake, err := s.conn.ReceiveHandshake()
	if err != nil {
		return err
	}

	if err := remoteHandshake.Validate(); err != nil {
		return err
	}

	handshake, err := protocol.NewHandshake(s.localID)
	if err != nil {
		return err
	}

	if err := s.conn.SendHandshake(handshake); err != nil {
		return err
	}

	s.remoteID = remoteHandshake.NodeID
	s.state = StateEstablished

	return nil
}

// Send sends a protocol message over the established session.
// Send is safe to call concurrently with Receive.
func (s *Session) Send(message protocol.Message) error {
	s.stateMu.RLock()
	state := s.state
	s.stateMu.RUnlock()

	if state != StateEstablished {
		return ErrNotEstablished
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	return s.conn.Send(message)
}

// Receive reads a protocol message from the established session.
// Receive is safe to call concurrently with Send.
func (s *Session) Receive() (protocol.Message, error) {
	s.stateMu.RLock()
	state := s.state
	s.stateMu.RUnlock()

	if state != StateEstablished {
		return protocol.Message{}, ErrNotEstablished
	}

	return s.conn.Receive()
}

func (s *Session) Close() error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.state == StateClosed {
		return nil
	}

	err := s.conn.Close()
	s.state = StateClosed

	return err
}

