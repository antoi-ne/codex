package com

import (
	"errors"
	"net"
	"strings"

	"github.com/antoi-ne/codex/internal/keys"
	"github.com/antoi-ne/codex/internal/store"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	sshConfig *ssh.ServerConfig
	address   string
}

var bannerMsg = strings.Replace(`
  ___  _____  ____  ____  _  _ 
 / __)(  _  )(  _ \( ___)( \/ )
( (__  )(_)(  )(_) ))__)  )  ( 
 \___)(_____)(____/(____)(_/\_)

`, "\n", "\n\r", -1)

func NewServer(kp *keys.KeyPair, address string) (s *Server) {
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
			return &ssh.Permissions{}, nil
		},
		BannerCallback: func(conn ssh.ConnMetadata) string {
			return bannerMsg
		},
	}

	s.sshConfig.AddHostKey(kp.Priv)

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

		sConn, chans, reqs, err := ssh.NewServerConn(c, s.sshConfig)
		if err != nil {
			continue
		}

		go ssh.DiscardRequests(reqs)
		go handleServerConn(sConn, chans)
	}
}

func handleServerConn(sConn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	for newChan := range chans {

		switch newChan.ChannelType() {
		case "session":
			newChan.Reject(ssh.Prohibited, "codex does not provide shell access")
		case "codex-init":
			ch, reqs, err := newChan.Accept()
			if err != nil {
				continue
			}
			go ssh.DiscardRequests(reqs)
			ch.Write([]byte("hello world!"))
			ch.Close()
		default:
			newChan.Reject(ssh.UnknownChannelType, "unknown channel type")
		}
	}
}
