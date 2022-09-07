package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"coulon.dev/codex/internal/configs"
	"coulon.dev/codex/internal/store"
)

var (
	configPathFlag string
)

func init() {
	flag.StringVar(&configPathFlag, "config", "/etc/codexd.json", "Config file path")
}

func main() {
	flag.Parse()

	if len(flag.Args()) > 1 {
		log.Fatalln("Too many arguments")
	}

	cfg, err := configs.LoadServerSettings(configPathFlag)
	if err != nil {
		log.Fatalf("could not parse the config file (%s)", err)
	}

	store.InitDB(cfg.DbPath)

	switch flag.Arg(0) {
	case "new":
		pubkey := prompt("SSH public key:")
		u, err := store.NewUser([]byte(pubkey), time.Now())
		if err != nil {
			log.Fatalf("could not parse the SSH public key (%s)", err)
		}
		name := prompt(fmt.Sprintf("Username (%s):", u.Name))
		if name != "" {
			u.Name = name
		}
		expInput := strings.TrimSpace(prompt("User expiration day (format: YYYY/MM/DD):"))
		exp, err := time.Parse("2006/01/02", expInput)
		if err != nil {
			log.Fatalf("Invalid expiration date (%s)", err)
		}
		u.Expiration = exp
		principals := prompt("Valid principals (space-separated):")
		u.Principals = strings.Split(principals, " ")
		err = u.Save()
		if err != nil {
			log.Fatalf("could not save the new user (%s)", err)
		}
		fmt.Printf("User %s created successfully\n", u.Name)
	case "update":
		log.Fatalln("Not implemented")
	case "":
		log.Fatalln("No arguments given")
	default:
		log.Fatalln("Unknown subcommand")
	}
}

func prompt(label string) (out string) {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s ", label)
	out, err := r.ReadString('\n')
	if err != nil {
		log.Fatalf("could not read user input")
	}
	return strings.TrimRight(out, "\n")
}
