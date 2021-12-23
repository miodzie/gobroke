package config

// TODO: move to app package?

import (
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

var config *Config
var dir string

type Config struct {
	Currency  string              `yaml:"currency"`
	Watchlist []string            `yaml:"watchlist"`
	SMTP      EmailNotifierConfig `yaml:"smtp"`
}

type EmailNotifierConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (edc EmailNotifierConfig) Auth() smtp.Auth {
	return smtp.PlainAuth("", edc.Username, edc.Password, edc.Host)
}

type Recipient struct {
	Default string `yaml:"default"`
	Name    string `yaml:"name"`
	Email   string `yaml:"email"`
	SMS     string `yaml:"sms"`
}

// Conf returns the Config singleton instance.
func Conf() *Config {
	return config
}

func Dir() string {
	return dir
}

func CreateDefaults() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dir = filepath.Join(home, "/.config/tendies")
	_, err = os.Stat(dir)
	if os.IsExist(err) {
		log.Debug().Msg("Config directory already exists.")
		return nil
	}

	err = os.MkdirAll(dir, 0755)

	if err != nil {
		log.Err(err).Msg("Failed to create default config directory.")
		return err
	}

	return nil
}

func CopyDefaultConfig() {

}

// Load parses, returns, and sets the singleton Config.
func Load(filepath string) (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, err
	}

	filepath = strings.Replace(filepath, "~", home, 1)

	config, err := Parse(filepath)
	if err != nil {
		return &Config{}, err
	}

	return config, nil
}

// Parse parses the yaml configuration file,
// returning a Config, or error.
func Parse(file string) (*Config, error) {
	var config *Config

	contents, err := os.ReadFile(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
