package core

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
)

var (
	configPath string
)

type Config struct {
	ListenAddr string `json:"listen"`
	RemoteAddr string `json:"remote"`
}

func init() {
	home, _ := homedir.Dir()
	configFilename := "chat-config.json"
	if len(os.Args) == 2 {
		configFilename = os.Args[1]
	}
	configPath = path.Join(home, configFilename)
}

func (config *Config) SaveConfig() {
	configJson, _ := json.MarshalIndent(config, "", "	")
	err := ioutil.WriteFile(configPath, configJson, 0644)
	if err != nil {
		fmt.Errorf("Save the configuration %s failed: %s", configPath, err)
	}
	log.Printf("Save the configuration %s successful\n", configPath)
}

func (config *Config) ReadConfig() {
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		log.Printf("Read the configuration from %s \n", configPath)
		file, err := os.Open(configPath)
		if err != nil {
			log.Fatalf("Read the configuration %s failed:%s", configPath, err)
		}
		defer file.Close()

		err = json.NewDecoder(file).Decode(config)
		if err != nil {
			log.Fatalf("Invaild JSON format configuration:\n%s", file)
		}
	}
}
