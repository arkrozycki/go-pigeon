package main

import (
	// b64 "encoding/base64"
	"fmt"
	// "reflect"
	// "io/ioutil"
	// "github.com/jhump/protoreflect"
	"github.com/jhump/protoreflect/desc"
	// "github.com/jhump/protoreflect/dynamic"

	mts "github.com/arkrozycki/go-pigeon/protorepo/message-transfer"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

func JSONPayloadToProto(js []byte, conf *Config) error {
	// log.Info().Msg(string(js))
	var r desc.ImportResolver
	// r.RegisterImportPath("message_transfer_requested.proto", "message_transfer_requested.proto")
	fd, err := r.LoadMessageDescriptor("MessageTransferRequested")
	if err != nil {
		log.Fatal().Err(err)
	}
	// md := fd.GetMessageTypes()
	// md := fd.FindMessage("MessageTransferRequested") //.(*desc.MessageDescriptor)
	// msg := dynamic.NewMessage(md)
	fmt.Printf("%+v\n", fd)

	msg := &mts.MessageTransferRequested{}
	err = protojson.Unmarshal(js, msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	m, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Debug().Msgf("%+v", msg)

	conn, err := Connect(conf)
	defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()
	err = Publish(ch, &conf.Publisher.Exchange, "cmd.amz.mts.send", m)
	if err != nil {
		panic(err)
	}

	return nil
}

/*
b64 "encoding/base64"
	"fmt"
	"io/ioutil"

	mts "github.com/arkrozycki/go-pigeon/protorepo/message-transfer"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"

data, _ := ioutil.ReadFile("test.jpg")
	sEnc := b64.StdEncoding.EncodeToString(data)
	log.Info().Msg(fmt.Sprint(len(data)))
	jsonData := fmt.Sprintf(`{
		"clientKey":"%s",
		"messageKey":"%s",
		"payload": "%s",
		"payloadFilename": "%s",
		"payloadSize": "%s"
	}`,
		"amz",
		"evt.amz.mts.order_file_received",
		sEnc,
		"test_json_remote.jpg",
		fmt.Sprint(int32(len(data))),
	)

	msg := &mts.MessageTransferRequested{}

	err = protojson.Unmarshal([]byte(jsonData), msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	msg.Payload, err = b64.StdEncoding.DecodeString(string(msg.Payload[:]))
	// log.Info().Msgf("msg: %+v", msg)

	m, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	// data, _ := ioutil.ReadFile("./tmp/test.jpg")
	// msg := new(mts.MessageTransferRequested)
	// msg.ClientKey = "amz"
	// msg.MessageKey = "evt.amz.mts.order_file_received"
	// msg.Payload = data
	// msg.PayloadSize = int32(len(data))
	// msg.PayloadFilename = "test_remote_send.jpg"
	// protoMsg, err := proto.Marshal(msg)

	log.Info().
		Str("ClientKey", msg.ClientKey).
		Str("MessageKey", msg.MessageKey).
		Str("PayloadFilename", msg.PayloadFilename).
		Str("PayloadSize", fmt.Sprint(msg.PayloadSize)).
		Msg("MSG")

	// decodedBody := new(mts.MessageTransferRequested)
	// err = proto.Unmarshal(protoMsg, decodedBody)
	// log.Printf("Decoded Body \t", decodedBody, "\n\n", "Original \t", data)

	conn, err := Connect(&conf)
	defer conn.Close()
	ch, err := GetAMQPChannel(conn)
	defer ch.Close()
	err = Publish(ch, &conf.Publisher.Exchange, "cmd.amz.mts.send", m)
	if err != nil {
		panic(err)
	}

*/
