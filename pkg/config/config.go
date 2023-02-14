package config

import (
	"evf/pkg/jira"
	"os"
	"syscall"

	"github.com/go-yaml/yaml"
	"golang.org/x/term"
)

// Config represents all the required configuration options
// declared in the `config.yaml` file
type Config struct {
	Errata struct {
		URL          string `yaml:"url"`
		KerberosConf string `yaml:"kerberos-conf"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		Realm        string `yaml:"realm"`
	}
	Jira struct {
		URL               string `yaml:"url"`
		Token             string `yanl:"token"`
		jira.SearchParams `yaml:"params"`
	}
}

// LoadConfig reads the `config.yaml` file
// and decodes its content
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

// InputPassword asks user to input their Kerberos password
func (config *Config) InputPassword() error {
	print("Input your kerberos password.\nPassword:")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}
	println()
	config.Errata.Password = string(bytePassword)
	return nil
}
