package store

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"time"

	"go.etcd.io/bbolt"
	"golang.org/x/crypto/ssh"
	"pkg.coulon.dev/codex/internal/keys"
)

type User struct {
	PubKey     []byte
	Name       string
	Serial     uint
	Principals []string
	Expiration time.Time
}

func NewUser(key []byte, exp time.Time) (u *User, err error) {
	u = new(User)

	pub, com, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return nil, err
	}

	u.PubKey = pub.Marshal()
	u.Name = com
	u.Serial = 0
	u.Principals = []string{}
	u.Expiration = exp

	return
}

func GetUser(pubkey []byte) (u *User, err error) {
	db, err := bbolt.Open(dbPath, 0666, &bbolt.Options{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var v []byte = nil

	if err = db.View(func(tx *bbolt.Tx) error {
		bx := tx.Bucket([]byte("users"))
		v = bx.Get(pubkey)
		return nil
	}); err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	b := bytes.NewBuffer(v)
	d := gob.NewDecoder(b)
	if err = d.Decode(&u); err != nil {
		return nil, err
	}

	return
}

func (u *User) Save() (err error) {
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return
	}
	defer db.Close()

	b := &bytes.Buffer{}
	e := gob.NewEncoder(b)
	if err = e.Encode(u); err != nil {
		return
	}

	if err = db.Update(func(tx *bbolt.Tx) error {
		bx, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return bx.Put(u.PubKey, b.Bytes())
	}); err != nil {
		return
	}

	return
}

func (u *User) MakeCertificate(ca *keys.KeyPair) (c *ssh.Certificate, err error) {
	pk, err := ssh.ParsePublicKey(u.PubKey)
	if err != nil {
		return nil, err
	}

	c = &ssh.Certificate{
		Key:             pk,
		Serial:          uint64(u.Serial) + 1,
		CertType:        ssh.UserCert,
		KeyId:           u.Name,
		ValidPrincipals: u.Principals,
		ValidAfter:      uint64(time.Now().Unix()),
		ValidBefore:     uint64(time.Now().Add(time.Hour * 24).Unix()),
	}

	if c.SignCert(rand.Reader, ca.Priv); err != nil {
		return nil, err
	}

	u.Serial += 1

	if err = u.Save(); err != nil {
		return nil, err
	}

	return
}
