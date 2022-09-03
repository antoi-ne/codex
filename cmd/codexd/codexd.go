package main

import (
	"flag"
	"log"

	"github.com/antoi-ne/codex/internal/com"
	"github.com/antoi-ne/codex/internal/keys"
	"github.com/antoi-ne/codex/internal/store"
)

var (
	addressFlag    string
	dbPathFlag     string
	caPubPathFlag  string
	caPrivPathFlag string
)

func init() {
	flag.StringVar(&addressFlag, "address", "127.0.0.1:2222", "Server address in the format ip:port")
	flag.StringVar(&dbPathFlag, "database", "./codex.db", "Path to the database file")
	flag.StringVar(&caPubPathFlag, "ca-pubkey", "./ca.pub", "Path to the CA's public key")
	flag.StringVar(&caPrivPathFlag, "ca-privkey", "./ca.key", "Path to the CA's private key")
}

func main() {
	flag.Parse()

	store.SetDbPath(dbPathFlag)

	kp, err := keys.LoadKeyPair(caPubPathFlag, caPrivPathFlag)
	if err != nil {
		log.Fatalf("could not parse the ca keys (%s)", err)
	}

	s := com.NewServer(kp, addressFlag)

	if err = s.Listen(); err != nil {
		log.Fatalf("error while running ssh server (%s)", err)
	}
}
