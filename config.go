package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Connection ConnectionConfig
	Publisher  PublisherConfig
	Consumer   ConsumerConfig
}

type ConnectionConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type PublisherConfig struct {
	Exchange ExchangeConfig
}

type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
}

type ConsumerConfig struct {
	Exchange ExchangeConfig
	Queue    string
	Bindings []string
}

// GetConf will retrieve the configuration file at the location (f).
func GetConf(f string) (Config, error) {
	var conf Config
	viper.SetConfigFile(f)
	// read the config
	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	// unmarshall into viper configuration object
	if err := viper.Unmarshal(&conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
