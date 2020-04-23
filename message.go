package main

import (
	"errors"
	"strconv"

	// ------------------------------------------
	// 		IMPORT YOUR PROTOBUF BINDINGS HERE
	// ------------------------------------------
	mts "github.com/arkrozycki/go-pigeon/protorepo/message-transfer"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// getProtoMessageByName
func getProtoMessageByName(name string) (proto.Message, error) {
	// ------------------------------------------
	// 		MAP YOUR MESSAGE TYPE TO PROTOBUF BINDING HERE
	// ------------------------------------------
	switch name {
	case "mts.MessageTransferRequested":
		return &mts.MessageTransferRequested{}, nil
	default:
		return nil, errors.New("Not found.")
	}

}

// Emit
func Emit(routingKey string, protoMessageName string, js []byte, conf *Config) error {
	var msg protoreflect.ProtoMessage
	msg, err := getProtoMessageByName(protoMessageName)

	// Unmarshal the json to protomessage
	err = protojson.Unmarshal(js, msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	// marshal the protomessage to the wire format
	m, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	// setup rabbit
	conn, err := Connect(conf)
	defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()

	log.Info().
		Str("routingKey", routingKey).
		Str("proto", protoMessageName).
		Str("msg", string(js)).
		Msg("PIGEON SENT")

	// publish message
	err = Publish(ch, &conf.Publisher.Exchange, routingKey, protoMessageName, m)
	if err != nil {
		panic(err)
	}

	return nil
}

// Consume
func Consume(conf *Config) error {
	// setup rabbit
	conn, err := Connect(conf)
	// defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	// defer ch.Close()

	q, err := ch.QueueDeclare(
		conf.Consumer.Queue.Name, // name
		false,                    // durable
		true,                     // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return err
	}

	for _, binding := range conf.Consumer.Queue.Bindings {
		err = ch.QueueBind(
			conf.Consumer.Queue.Name,    // queue name
			binding,                     // routing key
			conf.Consumer.Exchange.Name, // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	// read messages from queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		// loop forever listening on channel
		for msg := range msgs {
			var header string
			if msg.Headers["proto"] != nil {
				header = msg.Headers["proto"].(string)
			}
			log.Info().
				Str("Exchange", msg.Exchange).
				Str("Redelivered", strconv.FormatBool(msg.Redelivered)).
				Str("DeliveryTag", strconv.FormatUint(msg.DeliveryTag, 10)).
				Str("RoutingKey", msg.RoutingKey).
				Str("HeaderProto", header).
				Str("Body", string(msg.Body)).
				Msg("PIGEON ARRIVED")
		}
	}()

	return nil
}
