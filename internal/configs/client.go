package configs

import (
	"bytes"
	"encoding/json"
	"os"
)

type ClientSettings struct {
	ServerAddress string `json:"serverAddress"`
	PubKeyPath    string `json:"pubKeyPath"`
	PrivKeyPath   string `json:"privKeyPath"`
}

func LoadClientSettings(path string) (ss *ClientSettings, err error) {
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
