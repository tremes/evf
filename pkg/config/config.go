package config

import (
	"evf/pkg/bugzilla"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Bugzilla struct {
		URL                string `yaml:"url"`
		Key                string `yaml:"key"`
		bugzilla.BugParams `yaml:"params"`
	}
	Errata struct {
		URL          string `yaml:"url"`
		KerberosConf string `yaml:"kerberos-conf"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Realm        string `yaml:"realm"`
	}
}

func LoadConfig() (*Config, error) {
	configFile, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	var config Config
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
