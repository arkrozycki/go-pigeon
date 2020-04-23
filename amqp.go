package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

var AMQPConnURLTemplate = "amqp://%s:%s@%s:%s/"

// Connect
func Connect(c *Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(getAMQPConnURL(c))
	return conn, err
}

// GetAMQPChannel
func GetAMQPChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	return channel, err
}

// getAMQPConnURL
func getAMQPConnURL(c *Config) string {
	return fmt.Sprintf(AMQPConnURLTemplate, c.Connection.User, c.Connection.Pass, c.Connection.Host, c.Connection.Port)
}

// CheckExchangeExists
func CheckExchangeExists(ch *amqp.Channel, exchange ExchangeConfig) error {
	err := ch.ExchangeDeclarePassive(
		exchange.Name,       // name
		exchange.Type,       // kind
		exchange.Durable,    // durable
		exchange.AutoDelete, // autoDelete
		false,               //internal
		false,               // noWait
		nil,                 // args
	)

	// (name, kind string, durable, autoDelete, internal, noWait bool, args Table)
	return err
}

func Publish(ch *amqp.Channel, e *ExchangeConfig, routingKey string, msg []byte) error {
	err := ch.Publish(
		e.Name,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return err
}
