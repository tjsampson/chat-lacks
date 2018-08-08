package core

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// AppConfig Struct holds the application configuration values
type AppConfig struct {
	TCPPort   string
	HTTPPort  string
	HostAddr  string
	LogOutput string
}

var (
	// AppConfiguration Singleton Pointer Intance of AppConfig
	AppConfiguration *AppConfig
)

func defaultConfig() *AppConfig {
	log.Println(" --- falling back to default config --- ")
	return &AppConfig{
		TCPPort:   "3000",
		HTTPPort:  "4000",
		HostAddr:  "localhost",
		LogOutput: "lacks.log",
	}
}

// Returns the Application Config
// Notice Multiple search paths
func getConf() *AppConfig {
	log.Println("+++ resolve the application config +++")

	viper.SetConfigName("config")
	// configuration search path locations
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.chat-lacks")
	viper.AddConfigPath("/etc/chat-lacks/")

	err := viper.ReadInConfig()

	if err != nil {
		log.Println(fmt.Errorf("error reading configuration: %s ", err))
		return defaultConfig()
	}

	conf := &AppConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.Println(fmt.Errorf("error unmarshaling configuration: %s ", err))
		return defaultConfig()
	}

	return conf
}
func init() {
	AppConfiguration = getConf()
	log.Println(fmt.Sprintf("Host: %v", AppConfiguration.HostAddr))
	log.Println(fmt.Sprintf("Tcp: %v", AppConfiguration.TCPPort))
	log.Println(fmt.Sprintf("Http: %v", AppConfiguration.HTTPPort))
	log.Println(fmt.Sprintf("LogOutput: %v", AppConfiguration.LogOutput))
}
