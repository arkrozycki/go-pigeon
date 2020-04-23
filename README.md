# go-pigeon
Helper for publishing messages using protobuf payloads


```json
{
  "proto": "name of the ProtoMessage to invoke",
  "routingKey": "rabbitmq-routing-key",
  "msg": {
    "foo": "bar"
  }
}
```