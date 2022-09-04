package configs

import (
	"bytes"
	"encoding/json"
	"os"
	"time"
)

type ServerSettings struct {
	Address       string        `json:"address"`
	DbPath        string        `json:"dbPath"`
	CaPubicPath   string        `json:"caPublicPath"`
	CaPrivatePath string        `json:"caPrivatePath"`
	CertValidity  time.Duration `json:"certValidity"`
}

func LoadServerSettings(path string) (ss *ServerSettings, err error) {
	ss = new(ServerSettings)

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(raw)
	d := json.NewDecoder(b)
	if err = d.Decode(ss); err != nil {
		return nil, err
	}

	return
}
