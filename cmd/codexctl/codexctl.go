package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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

	store.InitDB(cfg.DbPath)

	key, err := os.ReadFile("./user.pub")
	if err != nil {
		log.Fatalf("could not open the public key (%s)", err)
	}

	u, err := store.NewUser(key, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatalf("could not create the user (%s)", err)
	}

	if err = u.Save(); err != nil {
		log.Fatalf("could not save the user in the database (%s)", err)
	}

	fmt.Println("user created successfully")
}
