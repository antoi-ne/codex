package main

import (
	"flag"
	"log"

	"github.com/antoi-ne/codex/internal/com"
	"github.com/antoi-ne/codex/internal/configs"
	"github.com/antoi-ne/codex/internal/keys"
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

	ca, err := keys.LoadKeyPair(cfg.CaPubicPath, cfg.CaPrivatePath)
	if err != nil {
		log.Fatalf("could not parse the ca keys (%s)", err)
	}

	s := com.NewServer(ca, cfg.Address)

	if err = s.Listen(); err != nil {
		log.Fatalf("error while running ssh server (%s)", err)
	}
}
