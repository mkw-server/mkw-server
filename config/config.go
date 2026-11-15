package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type ConfigData struct {
	ReplayFileOutPath      string
	GroupReplayFilesByDate bool
}

const ConfigFileName = "config.toml"

var config ConfigData

func LoadConfig() error {
	// Since wfc-server launches mkw-server, we look for the file relative to the mkw-server
	// executable path.
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return err
	}

	configPath := filepath.Join(filepath.Dir(exePath), ConfigFileName)
	_, err = toml.DecodeFile(configPath, &config)
	return err
}

func ReplayFileOutPath() string {
	return config.ReplayFileOutPath
}

func GroupReplayFilesByDate() bool {
	return config.GroupReplayFilesByDate
}
