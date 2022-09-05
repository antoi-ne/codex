package com

import (
	"errors"
	"html/template"
	"log"
	"net"
	"strings"

	"github.com/antoi-ne/codex/internal/keys"
	"github.com/antoi-ne/codex/internal/store"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	sshConfig *ssh.ServerConfig
	address   string
	ca        *keys.KeyPair
}

var bannerMsg = strings.Replace(`
  ___  _____  ____  ____  _  _ 
 / __)(  _  )(  _ \( ___)( \/ )
( (__  )(_)(  )(_) ))__)  )  ( 
 \___)(_____)(____/(____)(_/\_)

`, "\n", "\n\r", -1)

var successTmpl = template.Must(template.New("successTmpl").Parse(strings.Replace(`
Hello {{ .User }},
Here is your new certificate:

{{ .Cert }}

`, "\n", "\n\r", -1)))

func NewServer(ca *keys.KeyPair, address string) (s *Server) {
	s = new(Server)

	s.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			u, err := store.GetUser(key.Marshal())
			if err != nil {
				return nil, err
			}
			if u == nil {
				return nil, errors.New("user not found")
			}
			return &ssh.Permissions{
				Extensions: map[string]string{
					"user-pubkey": string(u.PubKey),
				},
			}, nil
		},
		BannerCallback: func(conn ssh.ConnMetadata) string {
			return bannerMsg
		},
	}

	s.address = address

	s.ca = ca

	s.sshConfig.AddHostKey(ca.Priv)

	return
}

func (s *Server) Listen() (err error) {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}

		go s.HandleConn(c)
	}
}

func (s *Server) HandleConn(nConn net.Conn) {
	conn, chans, reqs, err := ssh.NewServerConn(nConn, s.sshConfig)
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)

	for nChan := range chans {
		switch nChan.ChannelType() {
		case "session":
			ch, rqs, err := nChan.Accept()
			if err != nil {
				return
			}
			defer ch.Close()

			go func(in <-chan *ssh.Request) {
				ok := false
				for req := range in {
					switch req.Type {
					case "shell":
						fallthrough
					case "pty-req":
						ok = true
					}
					if req.WantReply {
						req.Reply(ok, nil)
					}
				}
			}(rqs)

			u, err := store.GetUser([]byte(conn.Permissions.Extensions["user-pubkey"]))
			if err != nil {
				return
			}
			if u == nil {
				panic("user sent from public key callback not found")
			}

			c, err := u.MakeCertificate(s.ca)
			if err != nil {
				log.Printf("could not generate certificate (%s)", err)
				ch.Close()
				return
			}

			if err = successTmpl.Execute(ch, struct {
				User string
				Cert string
			}{User: u.Name, Cert: string(ssh.MarshalAuthorizedKey(c))}); err != nil {
				log.Printf("could not execute template (%s)", err)
			}

			return

		default:
			nChan.Reject(ssh.UnknownChannelType, "unknown channel type")
		}
	}
}
