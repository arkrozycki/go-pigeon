package main

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"sync"

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
		// Str("msg", string(js)).
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
	defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()

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

	// loop over bindings configs and bind all routing keys to queue
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done() // basically never
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
				// Str("Body", string(msg.Body)).
				Msg("PIGEON ARRIVED")

			// no header, no marshal
			if header == "" {
				continue
			}

			var p protoreflect.ProtoMessage
			p, err := getProtoMessageByName(header)
			if err != nil {
				log.Err(err)
				continue
			}
			// wire to proto
			err = proto.Unmarshal(msg.Body, p)
			if err != nil {
				log.Err(err)
				continue
			}

			// proto to json
			json, err := protojson.Marshal(p)
			if err != nil {
				log.Err(err)
				continue
			}
			// log.Info().
			// 	Str("json", string(json)).
			// 	Msg("Body")

			// if we dont have a webhook config ignore the rest
			if conf.Consumer.Webhook.Uri == "" {
				continue
			}

			// send json to webhook
			r := bytes.NewReader(json)
			req, err := http.NewRequest(conf.Consumer.Webhook.Verb, conf.Consumer.Webhook.Uri, r)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Err(err)
				continue
			}
			bod := &bytes.Buffer{}
			_, err = bod.ReadFrom(resp.Body)
			if err != nil {
				log.Err(err)
				continue
			}
			resp.Body.Close()

			log.Info().
				Str("uri", conf.Consumer.Webhook.Uri).
				Str("verb", conf.Consumer.Webhook.Verb).
				Str("statusCode", strconv.Itoa(resp.StatusCode)).
				Msg("WEBHOOK")

			if resp.StatusCode != 200 && resp.StatusCode != 201 {
				log.Info().Str("err", "Remote did not return 200 or 201.").Msg("WEBHOOK")
				continue
			}

		}
	}()

	wg.Wait()
	return nil
}
