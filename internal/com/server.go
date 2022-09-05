package com

import (
	"errors"
	"html/template"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/antoi-ne/codex/internal/keys"
	"github.com/antoi-ne/codex/internal/store"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	sshConfig *ssh.ServerConfig
	address   string
	ca        *keys.KeyPair
	users     map[string]*store.User
	mu        sync.RWMutex
}

var bannerMsg = strings.Replace(`
  ___  _____  ____  ____  _  _ 
 / __)(  _  )(  _ \( ___)( \/ )
( (__  )(_)(  )(_) ))__)  )  ( 
 \___)(_____)(____/(____)(_/\_)
                         v0.0.0
`, "\n", "\n\r", -1)

var successTmpl = template.Must(template.New("successTmpl").Parse(strings.Replace(`
Hello {{ .User }},
Here is your new certificate:

{{ .Cert }}
Certificate details:

    Serial number: {{ .Serial }}
    Expiration date: {{ .Exp }}

`, "\n", "\n\r", -1)))

func NewServer(ca *keys.KeyPair, address string) (s *Server) {
	s = new(Server)

	s.sshConfig = &ssh.ServerConfig{
		PublicKeyCallback: s.publicKeyCallback,
		BannerCallback: func(conn ssh.ConnMetadata) string {
			return bannerMsg
		},
	}

	s.address = address

	s.ca = ca

	s.sshConfig.AddHostKey(ca.Priv)

	s.users = make(map[string]*store.User)

	return
}

func (s *Server) publicKeyCallback(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	u, err := store.GetUser(key.Marshal())
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("user not found")
	}

	s.mu.Lock()
	s.users[string(conn.SessionID())] = u
	s.mu.Unlock()

	return &ssh.Permissions{}, nil
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

		go s.handleConn(c)
	}
}

func (s *Server) handleConn(nConn net.Conn) {
	conn, chans, reqs, err := ssh.NewServerConn(nConn, s.sshConfig)
	if err != nil {
		return
	}

	defer func() {
		delete(s.users, string(conn.SessionID()))
	}()

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

			s.mu.RLock()
			u := s.users[string(conn.SessionID())]
			s.mu.RUnlock()

			if u == nil {
				return
			}

			c, err := u.MakeCertificate(s.ca)
			if err != nil {
				log.Printf("could not generate certificate (%s)", err)
				return
			}

			if err = successTmpl.Execute(ch, struct {
				User   string
				Cert   string
				Serial uint64
				Exp    string
			}{
				User:   u.Name,
				Cert:   "\033[35m" + string(ssh.MarshalAuthorizedKey(c)) + "\033[0m",
				Serial: uint64(u.Serial),
				Exp:    time.Unix(int64(c.ValidBefore), 0).Format(time.RFC822),
			}); err != nil {
				log.Printf("could not execute template (%s)", err)
			}

			return

		default:
			nChan.Reject(ssh.UnknownChannelType, "unknown channel type")
		}
	}
}
