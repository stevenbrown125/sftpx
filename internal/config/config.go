package config

import (
	"encoding/json"
	"os"
)

type SFTPConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	User           string `json:"user"`
	Password       string `json:"password,omitempty"`
	PrivateKeyPath string `json:"privateKeyPath,omitempty"`
	Passphrase     string `json:"passphrase,omitempty"`
}

type Config struct {
	WatchDir     string     `json:"watchDir"`
	RemoteDir    string     `json:"remoteDir"`
	LogDir       string     `json:"logDir"`
	LogFile      string     `json:"logFile"`
	DelaySeconds int        `json:"delaySeconds"`
	Workers      int        `json:"workers"`
	SFTP         SFTPConfig `json:"sftp"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
