# go-pigeon

Useful for testing message based microservices transported via RabbitMQ and serialized by protocol buffers.

# Who might this help?
If you need the ability to publish JSON messages via REST have them converted to protobuf message and then published to rabbitmq exchange.

If you need the ability to attach a queue to an exchange with specific message bindings and have the incoming messages converted from protobuf to JSON and sent to a webhook.

## Installation

Just run docker-compose up and you should have service running and listening on port configured in `config.yml` by default `:8080`.

```bash
docker-compose up
```

## Configuration

Configure your environment within the `config.yml` file.

This package doesn't come with any protobuf schemas or bindings. You will have to provide your own. To do that you will need to:

- edit the `message.go` source file
- import your protobuf bindings (e.g. `myProtos "github.com/arkrozycki/protos")`
- map your messages to the protobuf binding in the function `getProtoMessageByName`

## API Message Specification

- `proto` The protobuf message name (used for binding)
- `routingKey` The routing_key for RabbitMQ
- `msg` The JSON representation of the proto message. This object will be unmarshalled to a ProtoMessage and sent to RabbitMQ.

```json
{
  "proto": "name of the ProtoMessage to invoke",
  "routingKey": "rabbitmq-routing-key",
  "msg": {
    "foo": "bar"
  }
}
```

### Binary Data

What if you want to include binary data? Just Base64 the binary file, remove any line breaks.

```bash
openssl base64 < path/to/file.xml | tr -d '\n' | pbcopy # copies base64 string to clipboard
stat -f%z file.xml # file size
```

## Webhook

You can configure a consumer to deliver messages to a HTTP endpoint. Under the `config.yml` section consumer > webhook. You can use services like requestbin or even turn up a local httpbin.
