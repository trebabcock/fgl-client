package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config holds server configuration information
type Config struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// GetConfig loads the server configuration from the file
func GetConfig() *Config {
	config := Config{}
	configFile, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println("error opening config.json:", err)
	}
	byteValue, _ := ioutil.ReadAll(configFile)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		fmt.Println("error unmarshalling config.json:", err)
	}
	defer configFile.Close()

	return &config
}
