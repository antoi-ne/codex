package com

import (
	"coulon.dev/codex/internal/keys"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sshConfig *ssh.ClientConfig
	address   string
}

func NewClient(kp *keys.KeyPair, address string) (c *Client) {
	c = new(Client)

	c.sshConfig = &ssh.ClientConfig{
		User: kp.Comment,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(kp.Priv),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	c.address = address

	return
}

func (c *Client) Connect() (err error) {
	client, err := ssh.Dial("tcp", c.address, c.sshConfig)
	if err != nil {
		return
	}

	_, reqs, err := client.OpenChannel("codex-init", nil)
	if err != nil {
		return
	}

	go ssh.DiscardRequests(reqs)

	return
}
