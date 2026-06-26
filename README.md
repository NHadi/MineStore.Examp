# ministore

You're building ministore, a lightweight document store where documents can have a parent (forming a tree). Implement the storage interface and SQLite schema, wire up a CLI binary, and make it compile for both Linux/amd64 and darwin/arm64. No actual DB logic needed, stub implementations that satisfy the interface and let the tests pass. You are allowed to use AI Coding Assistant (chat-based, terminal-based, or AI tab completion), however, no usage of Coding Assistant adds more points.

## Constraints

- Do not modify `store/store_test.go`
- `make test` must pass
- `make build` must produce binaries for both `linux/amd64` and `darwin/arm64`
- (optional) Write short CLAUDE.md for this project
- (optional) Write agent skills for this project, like golang skills you usually use
- You have 30 minutes. If you finish this before 30 minutes, you can stop and start to explain the implementation, or continue adding features, test, infrastructure needed until times up before explaining the implementation.