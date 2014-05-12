package main

import (
	"github.com/streadway/amqp"
	"github.com/wurkhappy/WH-Config"
	"log"
)

type Event struct {
	Name string
	Body []byte
}

type Events []*Event

func (e *Event) Publish() {
	ch := getChannel()
	defer ch.Close()
	e.PublishOnChannel(ch)
}

func (e Events) Publish() {
	ch := getChannel()
	defer ch.Close()
	for _, event := range e {
		event.PublishOnChannel(ch)
	}
}

var connection *amqp.Connection

func dialRMQ() {
	var err error
	connection, err = amqp.Dial(config.RMQBroker)
	if err != nil {
		panic(err)
	}
}

func (e *Event) PublishOnChannel(ch *amqp.Channel) {
	if ch == nil {
		ch = getChannel()
		defer ch.Close()
	}

	ch.ExchangeDeclare(
		config.MainExchange, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)
	ch.Publish(
		config.MainExchange, // exchange
		e.Name,              // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        e.Body,
		})
}

func getChannel() *amqp.Channel {
	ch, err := connection.Channel()
	if err != nil {
		dialRMQ()
		ch, err = connection.Channel()
		if err != nil {
			log.Print(err.Error())
		}
	}

	return ch
}
