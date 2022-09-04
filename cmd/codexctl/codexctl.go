package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/antoi-ne/codex/internal/configs"
	"github.com/antoi-ne/codex/internal/store"
)

var (
	configPathFlag string
)

func init() {
	flag.StringVar(&configPathFlag, "config", "/etc/codexd.json", "Config file path")
}

func main() {
	flag.Parse()

	cfg, err := configs.LoadServerSettings(configPathFlag)
	if err != nil {
		log.Fatalf("could not parse the config file (%s)", err)
	}

	store.SetDbPath(cfg.DbPath)

	fmt.Println("not implemented")
}
