package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	// Search config in /etc, /opt, home directory, or current directory with name "config.yaml".
	viper.AddConfigPath("/etc/logstash-pipeline-tester/")
	viper.AddConfigPath("/opt/logstash-pipeline-tester/")
	viper.AddConfigPath("~/.logstash-pipeline-tester/")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
