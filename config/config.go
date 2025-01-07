package config

import (
	"github.com/jinzhu/configor"
)

type Configuration struct {
	Server struct {
		ListenAddr string `default:""`
		Port       string `default:"8080"`
	}
	Database struct {
		SQL struct {
			Connection string `default:"root:root@tcp(localhost:3306)/chatapp?charset=utf8mb4&parseTime=True&loc=Local"`
		}
		Redis struct {
			Addr string `default:"localhost:6379"`
		}
	}

	DefaultUser struct {
		Name     string `default:"admin"`
		Email    string `default:"admin@example.com"`
		Password string `default:"admin"`
	}
}

func configFiles() []string {
	return []string{"config.yml"}
}

// Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	err := configor.Load(conf, configFiles()...)
	if err != nil {
		panic(err)
	}
	return conf
}
