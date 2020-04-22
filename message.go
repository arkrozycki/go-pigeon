package main

import (
	// b64 "encoding/base64"

	// "reflect"
	// "io/ioutil"
	// "github.com/jhump/protoreflect"

	// "github.com/jhump/protoreflect/dynamic"

	"errors"

	mts "github.com/arkrozycki/go-pigeon/protorepo/message-transfer"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// getProtoMessageByName
func getProtoMessageByName(name string) (proto.Message, error) {
	switch name {
	case "MessageTransferRequested":
		return &mts.MessageTransferRequested{}, nil
	default:
		return nil, errors.New("Not found.")
	}

}

// JSONPayloadToProto
func JSONPayloadToProto(messageType string, js []byte, conf *Config) error {
	var msg protoreflect.ProtoMessage
	msg, err := getProtoMessageByName(messageType)

	err = protojson.Unmarshal(js, msg)
	if err != nil {
		log.Fatal().Err(err)
	}

	m, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal().Err(err)
	}

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
