package main

import (
	"flag"
	"log"

	"github.com/antoi-ne/codex/internal/com"
	"github.com/antoi-ne/codex/internal/keys"
)

var (
	addressFlag  string
	PubPathFlag  string
	PrivPathFlag string
)

func init() {
	flag.StringVar(&addressFlag, "address", "127.0.0.1:2222", "Server address in the format ip:port")
	flag.StringVar(&PubPathFlag, "pubkey", "./user.pub", "Path to the user's public key")
	flag.StringVar(&PrivPathFlag, "privkey", "./user.key", "Path to the user's private key")
}

func main() {
	flag.Parse()

	kp, err := keys.LoadKeyPair("./user.pub", "./user.key")
	if err != nil {
		log.Fatalf("could not parse the keys (%s)", err)

	}

	c := com.NewClient(kp, addressFlag)

	if err = c.Connect(); err != nil {
		log.Fatalf("could not connect to the server (%s)", err)
	}
}
