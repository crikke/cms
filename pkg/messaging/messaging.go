package messaging

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Shutdown()
}

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type ConnectionInfo struct {
	URI      string
	Exchange struct {
		Name       string
		Type       string
		Durable    bool
		AutoDelete bool
		Args       amqp.Table
	}
	Queue struct {
		Name       string
		Durable    bool
		AutoDelete bool
		Key        string
	}
}

func NewConsumer(info ConnectionInfo) error {

	c := &consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	var err error
	c.conn, err = amqp.Dial(info.URI)

	if err != nil {
		return err
	}

	go func() {
		fmt.Printf("%s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()

	if err != nil {
		return err
	}

	err = c.channel.ExchangeDeclare(
		info.Exchange.Name,
		info.Exchange.Type,
		info.Exchange.Durable,
		info.Exchange.AutoDelete,
		false,
		false,
		info.Exchange.Args)

	if err != nil {
		return err
	}

	// c.channel.QueueDeclare(
	// 	info.Queue.Name,
	// 	info.Queue.Durable,
	// 	info.Queue.AutoDelete,
	// )

	return nil
}
