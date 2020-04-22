package main

import (
	"io/ioutil"
	"os"

	mts "github.com/arkrozycki/go-pigeon/protorepo/message-transfer"
	"github.com/golang/protobuf/proto"
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
		panic(err) // launch must stop if we don't pass the checks
	}

	data, _ := ioutil.ReadFile("./771199.jpg")
	msg := new(mts.MessageTransferRequested)
	msg.ClientKey = "amz"
	msg.MessageKey = "evt.amz.mts.order_file_received"
	msg.Payload = data
	msg.PayloadSize = int32(len(data))
	msg.PayloadFilename = "771199_test_remote_send.jpg"
	protoMsg, err := proto.Marshal(msg)

	// decodedBody := new(mts.MessageTransferRequested)
	// err = proto.Unmarshal(protoMsg, decodedBody)
	// log.Printf("Decoded Body \t", decodedBody, "\n\n", "Original \t", data)

	conn, err := Connect(&conf)
	defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()
	err = Publish(ch, &conf.Publisher.Exchange, "cmd.amz.mts.send", protoMsg)
	if err != nil {
		panic(err)
	}
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
	log.Info().Msgf("OK ... AMQP Exchange %s exists", conf.Publisher.Exchange.Name)

	return nil
}
