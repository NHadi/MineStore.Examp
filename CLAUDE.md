# ministore — Project CLAUDE.md

A lightweight document store where documents optionally belong to a parent, forming a tree.
Used as a Go coding exercise for the Repository Engineer role.

## Tech Stack

- **Language**: Go 1.22+
- **Module**: `github.com/example/ministore`
- **Storage**: In-memory (`MemoryStore`) + optional SQLite (`SQLiteStore`, pure-Go via `modernc.org/sqlite`)
- **Build**: `make build` / `make cross`

## Go Installation

If `go` is not recognized, your Go toolchain is likely installed at a custom path.

**Windows (this machine):**
```
C:\Go\go\bin\go.exe
```
Add to PATH permanently:
```powershell
[Environment]::SetEnvironmentVariable("Path", "C:\Go\go\bin;$env:PATH", "User")
```

**Other common Go locations:**
- `C:\Program Files\Go\bin\go.exe`
- `C:\Go\bin\go.exe`
- `~/.local/go/bin/go` (Linux/macOS)

To find Go on your system:
```powershell
# Windows: search in common locations
Get-ChildItem -Path C:\ -Recurse -Filter "go.exe" -ErrorAction SilentlyContinue | Select-Object -First 5 FullName
```

## File Responsibilities

| File | Purpose |
| :--- | :--- |
| `store/store.go` | `Document` struct, `Store` interface, `MemoryStore` implementation, error sentinels (`ErrNotFound`, `ErrEmptyID`, `ErrParentNotFound`). |
| `store/sqlite.go` | `SQLiteStore` implementation. Build with `-tags sqlite`. |
| `cmd/ministore/main.go` | Memory-only CLI entry (no sqlite tag). |
| `cmd/ministore/sqlite_main.go` | Full CLI with SQLite support (`--db` flag). Build with `-tags sqlite`. |
| `Makefile` | `make build`, `make test`, `make cross`. |

## Quick Commands

Run these from the `ministrore-starter/` directory.

### Linux / macOS / Git Bash / WSL

```bash
# Run tests
go test -v ./...

# Build memory-only binary
go build -o bin/ministore ./cmd/ministore

# Build with SQLite support
go build -tags sqlite -o bin/ministore ./cmd/ministore

# Cross-compile
GOOS=linux GOARCH=amd64 go build -tags sqlite -o bin/ministore-linux-amd64 ./cmd/ministore
GOOS=darwin GOARCH=arm64 go build -tags sqlite -o bin/ministore-darwin-arm64 ./cmd/ministore

# Run CLI REPL
./bin/ministore
# Commands: put <id> <content>
#           get <id>
#           children <parentID>
#           exit
```

### Windows PowerShell / CMD

```powershell
# Navigate to project
cd C:\Users\nurul\Repositories\Excercise\ministrore-starter

# Run tests (use full path if 'go' not in PATH)
& "C:\Go\go\bin\go.exe" test -v ./...

# Build memory-only binary
& "C:\Go\go\bin\go.exe" build -o bin/ministore ./cmd/ministore

# Build with SQLite support
& "C:\Go\go\bin\go.exe" build -tags sqlite -o bin/ministore ./cmd/ministore

# Run CLI REPL
.\bin\ministore.exe
```

## Architecture

```
store.Store (interface)
  ├── store.MemoryStore   (in-memory map + RWMutex)
  └── store.SQLiteStore   (modernc.org/sqlite, build with -tags sqlite)
```

## Key Constraints

- **DO NOT modify `store/store_test.go`** — it is locked for scoring.
- All 3 tests must pass: `TestPutGet`, `TestChildren`, `TestMissingParent`.
- `store/store.go` must define `Document` with fields `ID string`, `ParentID *string`, `Content string`.
- `Store` interface methods: `Put(doc Document) error`, `Get(id string) (Document, error)`, `Children(parentID string) ([]Document, error)`.
