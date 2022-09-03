package main

import (
	"flag"
	"fmt"

	"github.com/antoi-ne/codex/internal/store"
)

var (
	dbPathFlag string
)

func init() {
	flag.StringVar(&dbPathFlag, "database", "./codex.db", "Path to the database file")
}

func main() {
	flag.Parse()

	store.SetDbPath(dbPathFlag)

	fmt.Println("not implemented")
}
