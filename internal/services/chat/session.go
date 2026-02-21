package chat

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ChatTurn represents a single turn in the conversation.
type ChatTurn struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"`
}

// Session holds all state for a single chat conversation.
type Session struct {
	ID                   string     `json:"id"`
	SelectedProvider     string     `json:"selected_provider"`      // "" = all providers
	LastMentionedProvider string    `json:"last_mentioned_provider"` // auto-detected from queries
	ChatHistory          []ChatTurn `json:"chat_history"`
	LastActivity         time.Time  `json:"last_activity"`
}

// LastNTurns returns the last n turns (individual messages, not pairs).
func (s *Session) LastNTurns(n int) []ChatTurn {
	if len(s.ChatHistory) <= n {
		return s.ChatHistory
	}
	return s.ChatHistory[len(s.ChatHistory)-n:]
}

// SessionStore manages all active chat sessions with TTL expiry.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	ttl      time.Duration
}

// NewSessionStore creates a new store and starts the cleanup goroutine.
func NewSessionStore(ttl time.Duration) *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*Session),
		ttl:      ttl,
	}
	go store.cleanupLoop()
	return store
}

// Create creates a new session and returns it.
func (s *SessionStore) Create() *Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &Session{
		ID:           uuid.New().String(),
		ChatHistory:  make([]ChatTurn, 0),
		LastActivity: time.Now(),
	}
	s.sessions[session.ID] = session
	return session
}

// Get retrieves a session by ID. Returns nil if not found or expired.
func (s *SessionStore) Get(id string) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil
	}
	if time.Since(session.LastActivity) > s.ttl {
		return nil
	}
	session.LastActivity = time.Now()
	return session
}

// Update saves changes to an existing session.
func (s *SessionStore) Update(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.LastActivity = time.Now()
	s.sessions[session.ID] = session
}

// cleanupLoop removes expired sessions every 5 minutes.
func (s *SessionStore) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for id, session := range s.sessions {
			if time.Since(session.LastActivity) > s.ttl {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
