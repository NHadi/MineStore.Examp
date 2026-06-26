package store

import (
	"errors"
	"sync"
)

// ErrNotFound is returned by Get when the requested document ID does not exist.
var ErrNotFound = errors.New("document not found")

// ErrEmptyID is returned when an ID argument is the empty string.
var ErrEmptyID = errors.New("id cannot be empty")

// ErrParentNotFound is returned when a ParentID references a non-existent document.
var ErrParentNotFound = errors.New("parent document not found")

// Document represents a stored document.
// ParentID is nil for root-level documents, and holds the parent's ID otherwise.
type Document struct {
	ID       string
	ParentID *string // nil means root-level
	Content  string
}

// Store defines the contract for document storage implementations.
type Store interface {
	// Put stores a document. If ParentID is non-nil, the referenced parent must
	// already exist in the store. Returns ErrParentNotFound if validation fails.
	Put(doc Document) error

	// Get retrieves a document by ID. Returns ErrNotFound if absent.
	Get(id string) (Document, error)

	// Children returns all documents whose ParentID equals parentID.
	// Returns an empty slice (not nil) if no children exist.
	Children(parentID string) ([]Document, error)
}

// MemoryStore is an in-memory, goroutine-safe implementation of Store.
type MemoryStore struct {
	mu  sync.RWMutex
	data map[string]Document
}

// NewMemoryStore returns an initialized MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]Document),
	}
}

// Put stores doc keyed by doc.ID. It validates ParentID if non-nil.
func (s *MemoryStore) Put(doc Document) error {
	if doc.ID == "" {
		return ErrEmptyID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if doc.ParentID != nil {
		if _, exists := s.data[*doc.ParentID]; !exists {
			return ErrParentNotFound
		}
	}

	s.data[doc.ID] = doc
	return nil
}

// Get returns the document with the given ID.
func (s *MemoryStore) Get(id string) (Document, error) {
	if id == "" {
		return Document{}, ErrEmptyID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, ok := s.data[id]
	if !ok {
		return Document{}, ErrNotFound
	}
	return doc, nil
}

// Children returns all documents whose ParentID matches parentID.
func (s *MemoryStore) Children(parentID string) ([]Document, error) {
	if parentID == "" {
		return nil, ErrEmptyID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Document
	for _, doc := range s.data {
		if doc.ParentID != nil && *doc.ParentID == parentID {
			result = append(result, doc)
		}
	}
	return result, nil
}
