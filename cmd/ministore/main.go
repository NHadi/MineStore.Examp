//go:build !sqlite

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/example/ministore/store"
)

const version = "v0.1.0"

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	showHelp    := flag.Bool("help",    false, "show help message")
	flag.Parse()

	if *showVersion {
		fmt.Printf("ministore %s\n", version)
		os.Exit(0)
	}
	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	s := store.NewMemoryStore()
	runREPL(os.Stdin, os.Stdout, os.Stderr, s)
}

func runREPL(rin *os.File, rout, rerr *os.File, s store.Store) {
	scanner := bufio.NewScanner(rin)
	fmt.Fprintln(rout, "ministrore REPL. Commands: put <id> [parentID] <content>, get <id>, children <parentID>, exit")

	for {
		fmt.Fprint(rout, "> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		cmd := strings.ToLower(args[0])

		var err error
		switch cmd {
		case "exit", "quit":
			return

		case "put":
			err = handlePut(rout, s, args[1:])

		case "get":
			err = handleGet(rout, s, args[1:])

		case "children":
			err = handleChildren(rout, s, args[1:])

		default:
			fmt.Fprintf(rerr, "unknown command: %q (try: put, get, children, exit)\n", cmd)
			continue
		}

		if err != nil {
			fmt.Fprintf(rerr, "error: %v\n", err)
		}
	}
}

func handlePut(w *os.File, s store.Store, args []string) error {
	if len(args) < 2 {
		return errors.New("usage: put <id> [parentID] <content>")
	}
	id := args[0]

	if len(args) == 2 {
		content := strings.Join(args[1:], " ")
		doc := store.Document{ID: id, Content: content}
		return s.Put(doc)
	}

	parentID := args[1]
	content := strings.Join(args[2:], " ")
	doc := store.Document{ID: id, ParentID: &parentID, Content: content}
	return s.Put(doc)
}

func handleGet(w *os.File, s store.Store, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: get <id>")
	}
	doc, err := s.Get(args[0])
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "ID: %s\n", doc.ID)
	if doc.ParentID != nil {
		fmt.Fprintf(w, "ParentID: %s\n", *doc.ParentID)
	}
	fmt.Fprintf(w, "Content: %s\n", doc.Content)
	return nil
}

func handleChildren(w *os.File, s store.Store, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: children <parentID>")
	}
	children, err := s.Children(args[0])
	if err != nil {
		return err
	}
	if len(children) == 0 {
		fmt.Fprintln(w, "(no children)")
		return nil
	}
	for _, child := range children {
		fmt.Fprintf(w, "  - ID: %s  Content: %s\n", child.ID, child.Content)
	}
	return nil
}

func printHelp() {
	fmt.Print(`ministrore - lightweight document store with tree structure

Usage:
  ministore [flags]

Flags:
  --version   print version and exit
  --help      show this help message
  --db <path> use SQLite store at <path> (omit for in-memory store)

REPL Commands:
  put <id> [parentID] <content>  store a document
  get <id>                       retrieve a document
  children <parentID>            list children of a document
  exit, quit                     exit the REPL
`)
}
