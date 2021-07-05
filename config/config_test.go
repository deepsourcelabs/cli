package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestCLIConfig_GetConfigPath(t *testing.T) {
	cfg := CLIConfig{}

	// No errors.
	configDirFn = func() (string, error) {
		return "/home/foo", nil
	}

	want := "/home/foo/.deepsource/config.toml"

	got, err := cfg.configPath()
	if got != want {
		t.Errorf("CLIConfig.GetConfigPath() = (%v,%v) want (%v,%v)", got, err, want, nil)
	}

	//Errors.
	configDirFn = func() (string, error) {
		return "", errors.New("random error")
	}

	want = ""

	got, err = cfg.configPath()
	if err == nil {
		t.Errorf("CLIConfig.GetConfigPath() = (%v,%v) want (%v,%v)", got, err, want, nil)
	}

}

func TestConfigManager_ReadFile(t *testing.T) {

	configDirFn = func() (string, error) {
		return "/home/foo", nil
	}

	want := CLIConfig{
		Token:          "thisistoken",
		RefreshToken:   "thisisrefreshtoken",
		TokenExpiresIn: 20,
	}

	//Mock os.readFile function.
	readFileFn = func(path string) ([]byte, error) {
		data, _ := toml.Marshal(want)
		return data, nil
	}

	got := CLIConfig{}
	err := got.ReadFile()
	if err != nil {
		t.Errorf("CLIConfig.ReadFile() = (%v, %v) want (%v, %v)", got, err, want, nil)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("CLIConfig.ReadFile() = (%v, %v) want (%v, %v)", got, err, want, nil)
	}

}

func TestCLIConfig_WriteFile(t *testing.T) {
	configDirFn = func() (string, error) {
		return os.TempDir(), nil
	}
	want := CLIConfig{
		Token:          "thisistoken",
		RefreshToken:   "thisisrefreshtoken",
		TokenExpiresIn: 20,
	}
	want.WriteFile()
	dir, _ := configDirFn()
	path := filepath.Join(dir, ConfigDirName, ConfigFileName)

	fmt.Println(path)
	_, err := os.Stat(path)
	if err != nil {
		t.Errorf("CLIConfig.WriteFile() = (%v)", err)
	}
}
