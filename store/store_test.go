package store_test

import (
	"testing"

	"github.com/example/ministore/store"
)

func TestPutGet(t *testing.T) {
	s := store.NewMemoryStore()

	doc := store.Document{
		ID:      "doc1",
		Content: "hello world",
	}
	if err := s.Put(doc); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got, err := s.Get("doc1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Content != doc.Content {
		t.Errorf("content: got %q, want %q", got.Content, doc.Content)
	}
}

func TestChildren(t *testing.T) {
	s := store.NewMemoryStore()

	parent := store.Document{ID: "parent", Content: "parent doc"}
	if err := s.Put(parent); err != nil {
		t.Fatalf("Put parent: %v", err)
	}

	parentID := "parent"
	child1 := store.Document{ID: "child1", ParentID: &parentID, Content: "child one"}
	child2 := store.Document{ID: "child2", ParentID: &parentID, Content: "child two"}

	if err := s.Put(child1); err != nil {
		t.Fatalf("Put child1: %v", err)
	}
	if err := s.Put(child2); err != nil {
		t.Fatalf("Put child2: %v", err)
	}

	children, err := s.Children("parent")
	if err != nil {
		t.Fatalf("Children: %v", err)
	}
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
}

func TestMissingParent(t *testing.T) {
	s := store.NewMemoryStore()

	parentID := "nonexistent"
	doc := store.Document{ID: "orphan", ParentID: &parentID, Content: "lost"}

	err := s.Put(doc)
	if err == nil {
		t.Fatal("expected error for missing parent, got nil")
	}
}
