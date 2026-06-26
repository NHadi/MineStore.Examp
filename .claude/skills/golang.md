# Go Agent Skill: ministore

## Context

`github.com/example/ministore` is a lightweight document store with a tree structure.
Used as a 30-minute coding exercise for Go candidates.

## Module Path

```
module github.com/example/ministore
```

## Adding a New Store Implementation

To add a new backend (e.g., PostgreSQL):

1. Create `store/pgstore.go` with:
   ```go
   package store

   type PGStore struct { ... }

   func NewPGStore(dsn string) (*PGStore, error) { ... }

   // Must implement all Store interface methods:
   // - Put(doc Document) error
   // - Get(id string) (Document, error)
   // - Children(parentID string) ([]Document, error)
   ```

2. Use `//go:build pg` tag if you want it separately compiled.
3. Update `cmd/ministore/main.go` (or create `pg_main.go`) to wire it up.

## Adding a CLI Subcommand

The CLI REPL is in `cmd/ministore/main.go` (memory) and `cmd/ministore/sqlite_main.go` (sqlite).

To add a new subcommand:

1. Add a case to the `switch cmd` block in `runREPL`:
   ```go
   case "mycommand":
       err = handleMyCommand(rout, s, args[1:])
   ```

2. Add the handler function:
   ```go
   func handleMyCommand(w *os.File, s store.Store, args []string) error {
       if len(args) < 1 {
           return errors.New("usage: mycommand <arg>")
       }
       // ...
   }
   ```

## Running Tests

```bash
# All tests (store package)
go test -v ./...

# With race detector (requires CGO)
go test -race -v ./...

# With SQLite tag
go test -tags sqlite -v ./...

# Both
go test -tags sqlite -race -v ./...
```

## Adding Dependencies

```bash
go get <package>
go mod tidy
```

SQLite driver (pure Go, no CGO):
```bash
go get modernc.org/sqlite
```

## Build Tags

| Tag | Description |
| :--- | :--- |
| none | Memory-only build (`MemoryStore`) |
| `sqlite` | Full build with SQLite support (`SQLiteStore`) |

Cross-compile without tags on the respective platform, or use:
```bash
GOOS=linux GOARCH=amd64 go build -tags sqlite -o bin/ministore ./cmd/ministore
GOOS=darwin GOARCH=arm64 go build -tags sqlite -o bin/ministore ./cmd/ministore
```

## Error Sentinels (exported from `store` package)

```go
var ErrNotFound      = errors.New("document not found")
var ErrEmptyID       = errors.New("id cannot be empty")
var ErrParentNotFound = errors.New("parent document not found")
```

## Interface Contract

```go
type Store interface {
    Put(doc Document) error
    Get(id string) (Document, error)
    Children(parentID string) ([]Document, error)
}

type Document struct {
    ID       string
    ParentID *string // nil = root-level
    Content  string
}
```

## Common Pitfalls

1. **Race on MemoryStore**: Always use `sync.RWMutex`. `Put` → `Lock()`, `Get`/`Children` → `RLock()`.
2. **Test file modified**: Never touch `store/store_test.go`.
3. **Empty `id`**: Both `Get` and `Children` must return `ErrEmptyID` for empty string.
4. **Nil slice vs empty slice**: `Children` returns `nil` on error but `[]Document{}` (empty non-nil) on success with zero results.
5. **Go version**: Target Go 1.22 in `go.mod`.
