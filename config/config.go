package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/Le-BlitzZz/real-time-chat-app/mode"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Server struct {
		ListenAddr string `yaml:"ListenAddr"`
		Port       string `yaml:"Port"`
	} `yaml:"Server"`
	Database struct {
		SQL struct {
			Connection string `yaml:"Connection"`
		} `yaml:"SQL"`
		Redis struct {
			Addr string `yaml:"Addr"`
		} `yaml:"Redis"`
	} `yaml:"Database"`
	DefaultUser struct {
		Name     string `yaml:"Name"`
		Email    string `yaml:"Email"`
		Password string `yaml:"Password"`
	} `yaml:"DefaultUser"`
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	err := conf.load()
	if err != nil {
		panic(err)
	}
	return conf
}

func getConfigFile() string {
	if mode.IsLocalDev() {
		return mode.ConfigLocalDevFile
	}
	return mode.ConfigDockerDevFile
}

func (conf *Configuration) load() error {
	configFile := getConfigFile()
	if !fileExists(configFile) {
		return errors.New(fmt.Sprintf("%s not found", configFile))
	}

	yamlConfig, err := os.ReadFile(configFile)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlConfig, conf)
}

func fileExists(fileName string) bool {
	if fileName == "" {
		return false
	}

	info, err := os.Stat(fileName)

	return err == nil && !info.IsDir()
}
