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

func LoadClientSettings(path string) (cs *ClientSettings, err error) {
	cs = new(ClientSettings)

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(raw)
	d := json.NewDecoder(b)
	if err = d.Decode(cs); err != nil {
		return nil, err
	}

	return
}
