package main

import (
	"flag"
	"log"

	"github.com/antoi-ne/codex/internal/com"
	"github.com/antoi-ne/codex/internal/configs"
	"github.com/antoi-ne/codex/internal/keys"
)

var (
	configPathFlag string
)

func init() {
	flag.StringVar(&configPathFlag, "config", "~/.codex.json", "Config file path")
}

func main() {
	flag.Parse()

	cfg, err := configs.LoadClientSettings(configPathFlag)
	if err != nil {
		log.Fatalf("could not parse the config file (%s)", err)
	}

	kp, err := keys.LoadKeyPair(cfg.PubKeyPath, cfg.PrivKeyPath)
	if err != nil {
		log.Fatalf("could not parse the keys (%s)", err)
	}

	c := com.NewClient(kp, cfg.ServerAddress)

	if err = c.Connect(); err != nil {
		log.Fatalf("could not connect to the server (%s)", err)
	}
}
