package store

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"time"

	"github.com/antoi-ne/codex/internal/keys"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/ssh"
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
	u.Principals = []string{"root"}
	u.Expiration = exp

	return
}

func GetUser(pubkey []byte) (u *User, err error) {
	db, err := bbolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var v []byte = nil

	if err = db.View(func(tx *bbolt.Tx) error {
		bx := tx.Bucket([]byte("users"))
		v = bx.Get([]byte(pubkey))
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
		err = bx.Put([]byte(u.PubKey), b.Bytes())
		return err
	}); err != nil {
		return
	}

	return
}

func (u *User) MakeCertificate(ca keys.KeyPair) (c *ssh.Certificate, err error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey(u.PubKey)
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
		ValidBefore:     minUint64(uint64(u.Expiration.Unix()), uint64(time.Now().Add(time.Hour*24).Unix())),
	}

	if c.SignCert(rand.Reader, ca.Priv); err != nil {
		return nil, err
	}

	u.Serial += 1

	u.Save()

	return
}

func minUint64(n1 uint64, n2 uint64) uint64 {
	if n1 < n2 {
		return n1
	}
	return n2
}
