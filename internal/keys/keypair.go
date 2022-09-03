package keys

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type KeyPair struct {
	Comment string
	Pub     ssh.PublicKey
	Priv    ssh.Signer
}

func LoadKeyPair(pubPath string, privPath string) (kp *KeyPair, err error) {
	kp = new(KeyPair)

	f, err := ioutil.ReadFile(pubPath)
	if err != nil {
		return nil, err
	}
	kp.Pub, kp.Comment, _, _, err = ssh.ParseAuthorizedKey(f)
	if err != nil {
		return nil, err
	}

	f, err = ioutil.ReadFile(privPath)
	if err != nil {
		return nil, err
	}
	kp.Priv, err = ssh.ParsePrivateKey(f)
	if err != nil {
		return nil, err
	}

	return
}
