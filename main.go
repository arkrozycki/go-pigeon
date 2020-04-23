package main

import (
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("STARTING UP")

	conf, err := GetConf("config.yml")
	if err != nil {
		panic(err) // no config, no start
	}

	log.Info().
		Str("host", conf.Connection.Host).
		Str("port", conf.Connection.Port).
		Str("user", conf.Connection.User).
		Msg("CONFIG")

	if err := launchStatusCheck(&conf); err != nil {
		panic(err) // we must stop if we don't pass the checks
	}

	// load up any consumer
	if err = Consume(&conf); err != nil {
		panic(err)
	}
	log.Info().
		Str("queue", conf.Consumer.Queue.Name).
		Str("exchange", conf.Consumer.Exchange.Name).
		Str("bindings", strings.Join(conf.Consumer.Queue.Bindings, " ")).
		Msg("OK ... AMQP Queue")

	// startup API listener
	var wg sync.WaitGroup
	wg.Add(1)
	go Serve(":8080", &conf, &wg)
	log.Info().
		Str("uri", "http://localhost:8080").
		Msg("OK ... REST API")
	wg.Wait()

}

// launchStatusCheck ensures we have everything we need to startup
func launchStatusCheck(conf *Config) error {
	// verify we can connect to message bus
	conn, err := Connect(conf)
	defer conn.Close()
	if err != nil {
		return err
	}
	log.Info().Msg("OK ... AMQP Connection")

	// verify we get a message channel
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()
	if err != nil {
		return err
	}
	log.Info().Msg("OK ... AMQP Channel")

	// verify the pre-configured exchange exists
	if err := CheckExchangeExists(ch, conf.Publisher.Exchange); err != nil {
		return err
	}
	log.Info().
		Str("exchange", conf.Publisher.Exchange.Name).
		Msg("OK ... AMQP Exchange")

	return nil
}
