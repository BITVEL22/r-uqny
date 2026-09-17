package peer

import (
	"errors"
	"time"

	"github.com/BITVEL22/r-uqny/internal/identity"
	"github.com/BITVEL22/r-uqny/internal/session"
)

var (
	ErrSessionAlreadyAttached = errors.New("session is already attached")
	ErrSessionNotAttached     = errors.New("session is not attached")
)

type State string

const (
	StateDisconnected State = "disconnected"
	StateConnecting   State = "connecting"
	StateHandshaking  State = "handshaking"
	StateConnected    State = "connected"
	StateClosed       State = "closed"
)

type Peer struct {
	ID            string
	Address       string
	State         State
	LastConnected time.Time
	LastSeen      time.Time
	Session       *session.Session
}

func New(id string, address string) (Peer, error) {
	if id == "" {
		return Peer{}, errors.New("peer ID cannot be empty")
	}

	if len(id) != identity.NodeIDLength {
		return Peer{}, errors.New("invalid peer ID")
	}

	if address == "" {
		return Peer{}, errors.New("peer address cannot be empty")
	}

	return Peer{
		ID:      id,
		Address: address,
		State:   StateDisconnected,
	}, nil
}

func (p *Peer) SetState(state State) error {
	switch state {
	case StateDisconnected,
		StateConnecting,
		StateHandshaking,
		StateConnected,
		StateClosed:
		p.State = state
		return nil
	default:
		return errors.New("invalid peer state")
	}
}

func (p *Peer) MarkConnected() {
	now := time.Now().UTC()

	p.State = StateConnected
	p.LastConnected = now
	p.LastSeen = now
}

func (p *Peer) MarkSeen() {
	p.LastSeen = time.Now().UTC()
}

func (p *Peer) AttachSession(s *session.Session) error {
	if s == nil {
		return errors.New("session cannot be nil")
	}

	if p.Session != nil {
		return ErrSessionAlreadyAttached
	}

	p.Session = s
	return nil
}

func (p *Peer) DetachSession() error {
	if p.Session == nil {
		return ErrSessionNotAttached
	}

	p.Session = nil
	return nil
}

func (p *Peer) HasSession() bool {
	return p.Session != nil
}

func (p *Peer) GetSession() (*session.Session, error) {
	if p.Session == nil {
		return nil, ErrSessionNotAttached
	}

	return p.Session, nil
}
