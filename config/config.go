package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml"
)

var configDirFn = os.UserHomeDir
var readFileFn = os.ReadFile
var createFn = os.Create

const (
	ConfigDirName  = "/.deepsource/"
	ConfigFileName = "/config.toml"
)

type CLIConfig struct {
	Token                 string
	RefreshToken          string
	TokenExpiresIn        time.Time
	RefreshTokenExpiresIn int64
}

type ConfigManager struct{}

func (cfg *CLIConfig) SetTokenExpiry(str string) {
	layout := "2006-01-02T15:04:05.999999999"
	tokenExpiresIn, _ := time.Parse(layout, str)
	cfg.TokenExpiresIn = tokenExpiresIn
}

func (cfg CLIConfig) IsExpired() bool {
	return time.Now().After(cfg.TokenExpiresIn)
}

//configDir returns the directory to store the config file.
func (CLIConfig) configDir() (string, error) {
	home, err := configDirFn()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDirName), nil
}

//confgPath returns the file path to the config file.
func (cfg CLIConfig) configPath() (string, error) {
	home, err := cfg.configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigFileName), nil
}

//ReadFile reads the CLI config file.
func (cfg *CLIConfig) ReadFile() error {
	path, err := cfg.configPath()
	if err != nil {
		return err
	}
	data, err := readFileFn(path)
	if err != nil {
		return err
	}
	return toml.Unmarshal(data, cfg)
}

//WriteFile writes the CLI config to file.
func (cfg *CLIConfig) WriteFile() error {

	data, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	configDir, err := cfg.configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return err
	}

	path, err := cfg.configPath()
	if err != nil {
		return err
	}

	// Create file
	file, err := os.Create(path)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.Write(data)

	return err
}

func (cfg *CLIConfig) Delete() error {
	path, err := cfg.configPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}
