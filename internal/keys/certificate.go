package keys

import (
	"crypto/rand"
	"time"

	"github.com/antoi-ne/codex/internal/store"
	"golang.org/x/crypto/ssh"
)

func minUint64(n1 uint64, n2 uint64) uint64 {
	if n1 < n2 {
		return n1
	}
	return n2
}

func (kp *KeyPair) NewCertificate(u *store.User) (cert *ssh.Certificate, err error) {
	cert = &ssh.Certificate{
		Key:         nil,
		Serial:      uint64(u.Serial),
		CertType:    ssh.UserCert,
		KeyId:       u.Name,
		ValidAfter:  uint64(time.Now().Unix()),
		ValidBefore: minUint64(uint64(u.Expiration.Unix()), uint64(time.Now().Add(time.Hour*24).Unix())),
		ValidPrincipals: []string{
			"root",
		},
	}

	if err = cert.SignCert(rand.Reader, kp.Priv); err != nil {
		return
	}

	return
}
