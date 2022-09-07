package main

import (
	"flag"
	"log"

	"pkg.coulon.dev/codex/internal/com"
	"pkg.coulon.dev/codex/internal/configs"
	"pkg.coulon.dev/codex/internal/keys"
	"pkg.coulon.dev/codex/internal/store"
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

	store.InitDB(cfg.DbPath)

	ca, err := keys.LoadKeyPair(cfg.CaPubicPath, cfg.CaPrivatePath)
	if err != nil {
		log.Fatalf("could not parse the ca keys (%s)", err)
	}

	s := com.NewServer(ca, cfg.Address)

	if err = s.Listen(); err != nil {
		log.Fatalf("error while running ssh server (%s)", err)
	}
}
