package main

import (
	"github.com/spf13/viper"
)

// Config
type Config struct {
	Connection ConnectionConfig
	Publisher  PublisherConfig
	Consumer   ConsumerConfig
}

// ConnectionConfig
type ConnectionConfig struct {
	Host string
	Port string
	User string
	Pass string
}

// PublisherConfig
type PublisherConfig struct {
	Exchange ExchangeConfig
}

// ConsumerConfig
type ConsumerConfig struct {
	Exchange ExchangeConfig
	Queue    QueueConfig
	Webhook  WebhookConfig
}

// QueueConfig
type QueueConfig struct {
	Name     string
	Bindings []string
}

// ExchangeConfig
type ExchangeConfig struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
}

// WebhookConfig
type WebhookConfig struct {
	Uri  string
	Verb string
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
