package config

import (
	vars "github.com/out-of-mind/catalog/variables"

	"encoding/json"
	"io/ioutil"
	"strings"
	"os"
)

type Config struct {
	DB_NAME string `json:"DB_NAME"`
	DB_USER string `json:"DB_USER"`
	DB_PASSOWRD string `json:"DB_PASSOWRD"`
	DB_SSLMODE string `json:"DB_SSLMODE"`

	REDIS_IP string `json:"REDIS_IP"`
	REDIS_PORT string `json:"REDIS_PORT"`
	REDIS_PASSWORD string `json:"REDIS_PASSWORD"`
	REDIS_DB int `json:"REDIS_DB"`

	TEMPLATES_DIR string `json:"TEMPLATES_DIR"`

	SECRET string `json:"SECRET"`
}

func ParseConfig(path string) Config {
	data, err := ioutil.ReadFile(path)
    if err != nil {
    	vars.Log.Fatal("Falied read config file, error: ", err)
    }

    var config Config

    json.Unmarshal(data, &config)

    if config.DB_NAME == "" {
		vars.Log.Fatal("DB_NAME cannot be empty")
	}
	if config.DB_USER == "" {
		vars.Log.Fatal("DB_USER cannot be empty")
	}
	if strings.Contains(config.DB_PASSOWRD, "ENV") {
		config.DB_PASSOWRD = os.Getenv("ENV_CATALOG_DB_PASSOWRD")
	}
	if config.DB_SSLMODE == "" {
		vars.Log.Println("DB_SSLMODE cannot be empty")
		os.Exit(1)
	}

	if config.REDIS_IP == "" {
		vars.Log.Fatal("REDIS_IP cannot be empty")
	}
	if config.REDIS_PORT == "" {
		vars.Log.Fatal("REDIS_PORT cannot be empty")
	}
	if strings.Contains(config.REDIS_PASSWORD, "ENV") {
		config.REDIS_PASSWORD = os.Getenv("ENV_CATALOG_REDIS_PASSWORD")
	}

	if config.TEMPLATES_DIR == "" {
		vars.Log.Fatal("TEMPLATES_DIR cannot be empty")
	} else {
		_, err = ioutil.ReadDir(config.TEMPLATES_DIR)
		if err != nil {	
			vars.Log.Fatal("Cannot read TEMPLATES_DIR: ", config.TEMPLATES_DIR)
		}
	}
	if config.SECRET == "" {
		vars.Log.Fatal("SECRET cannot be empty")
	} else if strings.Contains(config.SECRET, "ENV") {
		config.SECRET = os.Getenv("ENV_CATALOG_SECRET")
	}

	vars.Secret = []byte(config.SECRET)
	vars.TemplateDir = config.TEMPLATES_DIR

	return config
}