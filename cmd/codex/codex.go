package main

import (
	"flag"
	"log"

	"pkg.coulon.dev/codex/internal/com"
	"pkg.coulon.dev/codex/internal/configs"
	"pkg.coulon.dev/codex/internal/keys"
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
