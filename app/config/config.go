package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/indrenicloud/tricloud-server/app/logg"
)

var conf *Config
var confFile = "./.meta/data/config.json"

type Config struct {
	Dev            bool
	DBpath         string
	StatDBpath     string
	EventProviders map[string]*EventConfig
}

type EventConfig struct {
	Apikey       string
	ConfigFile   string
	TokenPerUser map[string][]string
	Options      map[string]string
}

func (c *Config) Update() {
	// invoker should lock
	rawc, err := json.Marshal(c)
	if err != nil {
		logg.Error("Could not save config:")
		logg.Error(err)
		os.Exit(1)
	}
	ioutil.WriteFile(confFile, rawc, 0644)
}

func init() {
	logg.Info("config init")
	data, err := ioutil.ReadFile(confFile)
	if err != nil {
		logg.Error("Could not read config:")
		logg.Error(err)
		os.Exit(1)
	}
	conf = &Config{}

	err = json.Unmarshal(data, conf)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}
}

func GetConfig() *Config {
	return conf
}
