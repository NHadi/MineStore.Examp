//go:build sqlite

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

// SQLiteStore is a SQLite-backed implementation of Store.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) a SQLite database at path and applies
// the documents table schema. It returns a ready-to-use store.
func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path+"?_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("store.NewSQLiteStore: sql.Open: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("store.NewSQLiteStore: Ping: %w", err)
	}

	if err := migrate(context.Background(), db); err != nil {
		return nil, fmt.Errorf("store.NewSQLiteStore: migrate: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

// migrate creates the documents table if it does not exist.
func migrate(ctx context.Context, db *sql.DB) error {
	const ddl = `
	CREATE TABLE IF NOT EXISTS documents (
		id         TEXT PRIMARY KEY,
		parent_id  TEXT,
		content    TEXT NOT NULL,
		FOREIGN KEY (parent_id) REFERENCES documents(id)
	);
	CREATE INDEX IF NOT EXISTS idx_documents_parent_id ON documents(parent_id);
	`
	_, err := db.ExecContext(ctx, ddl)
	return err
}

// Put inserts or replaces a document. ParentID is stored as NULL if nil.
func (s *SQLiteStore) Put(doc Document) error {
	if doc.ID == "" {
		return ErrEmptyID
	}

	ctx := context.Background()

	if doc.ParentID != nil {
		var exists bool
		err := s.db.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM documents WHERE id = ?)", *doc.ParentID,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("store.SQLiteStore.Put: parent check: %w", err)
		}
		if !exists {
			return ErrParentNotFound
		}
	}

	var parentID any = nil
	if doc.ParentID != nil {
		parentID = *doc.ParentID
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT OR REPLACE INTO documents (id, parent_id, content) VALUES (?, ?, ?)`,
		doc.ID, parentID, doc.Content,
	)
	if err != nil {
		return fmt.Errorf("store.SQLiteStore.Put: INSERT: %w", err)
	}
	return nil
}

// Get returns the document with the given ID.
func (s *SQLiteStore) Get(id string) (Document, error) {
	if id == "" {
		return Document{}, ErrEmptyID
	}

	ctx := context.Background()
	var doc Document
	var parentID sql.NullString

	err := s.db.QueryRowContext(ctx,
		"SELECT id, parent_id, content FROM documents WHERE id = ?", id,
	).Scan(&doc.ID, &parentID, &doc.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Document{}, ErrNotFound
		}
		return Document{}, fmt.Errorf("store.SQLiteStore.Get: %w", err)
	}

	if parentID.Valid {
		doc.ParentID = &parentID.String
	}

	return doc, nil
}

// Children returns all documents whose parent_id equals parentID.
func (s *SQLiteStore) Children(parentID string) ([]Document, error) {
	if parentID == "" {
		return nil, ErrEmptyID
	}

	ctx := context.Background()
	rows, err := s.db.QueryContext(ctx,
		"SELECT id, parent_id, content FROM documents WHERE parent_id = ?", parentID,
	)
	if err != nil {
		return nil, fmt.Errorf("store.SQLiteStore.Children: %w", err)
	}
	defer rows.Close()

	var children []Document
	for rows.Next() {
		var doc Document
		var pid sql.NullString
		if err := rows.Scan(&doc.ID, &pid, &doc.Content); err != nil {
			return nil, fmt.Errorf("store.SQLiteStore.Children: row scan: %w", err)
		}
		if pid.Valid {
			doc.ParentID = &pid.String
		}
		children = append(children, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store.SQLiteStore.Children: rows: %w", err)
	}
	return children, nil
}

// Close releases the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
