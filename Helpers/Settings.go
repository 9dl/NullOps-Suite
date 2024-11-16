package Helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type ConfigurationManager struct {
	configMu sync.RWMutex
	config   Configuration
}

type Configuration struct {
	PrintInvalid  bool   `json:"PrintInvalid"`
	Threads       int    `json:"Threads"`
	Timeout       int    `json:"Timeout"`
	ThreadingType string `json:"ThreadingType"`
}

func NewConfigurationManager() *ConfigurationManager {
	return &ConfigurationManager{}
}

func (cm *ConfigurationManager) LoadConfig() {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		cm.config = Configuration{PrintInvalid: false, Threads: 550, ThreadingType: determineStrategy(), Timeout: 5000}
		cm.saveConfig()
		return
	}

	err = json.Unmarshal(data, &cm.config)
	if err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}
}

func (cm *ConfigurationManager) GetPrintInvalid() bool {
	cm.configMu.RLock()
	defer cm.configMu.RUnlock()

	return cm.config.PrintInvalid
}

func (cm *ConfigurationManager) GetThreadingType() string {
	cm.configMu.RLock()
	defer cm.configMu.RUnlock()

	return cm.config.ThreadingType
}

func (cm *ConfigurationManager) GetThreads() int {
	cm.configMu.RLock()
	defer cm.configMu.RUnlock()

	return cm.config.Threads
}

func (cm *ConfigurationManager) GetTimeout() int {
	cm.configMu.RLock()
	defer cm.configMu.RUnlock()

	return cm.config.Timeout
}

func (cm *ConfigurationManager) SetPrintInvalid(value bool) {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	cm.config.PrintInvalid = value
	cm.saveConfig()
}

func (cm *ConfigurationManager) SetThreads(value int) {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	cm.config.Threads = value
	cm.saveConfig()
}

func (cm *ConfigurationManager) SetTimeout(value int) {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	cm.config.Timeout = value
	cm.saveConfig()
}

func (cm *ConfigurationManager) SetThreadingType(value string) {
	cm.configMu.Lock()
	defer cm.configMu.Unlock()

	cm.config.ThreadingType = value
	cm.saveConfig()
}

func (cm *ConfigurationManager) saveConfig() {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling config:", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("config.json", data, 0600)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		os.Exit(1)
	}
}
