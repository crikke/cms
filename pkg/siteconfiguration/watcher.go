package siteconfiguration

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/streadway/amqp"
)

const cfgExchange = "cms.siteconfiguration"

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

// Initializes a temporary queue that subscribes to configuration changes
func NewConfigurationWatcher(uri string, cfg *SiteConfiguration) (io.Closer, error) {
	c := &consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	var err error
	c.conn, err = amqp.Dial(uri)

	if err != nil {
		return nil, err
	}

	go func() {
		fmt.Printf("%s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	c.channel, err = c.conn.Channel()

	if err != nil {
		return nil, err
	}

	err = c.channel.ExchangeDeclare(
		cfgExchange,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		return nil, err
	}

	q, err := c.channel.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = c.channel.QueueBind(
		q.Name,
		"",
		cfgExchange,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	messages, err := c.channel.Consume(q.Name, "", false, true, false, false, nil)

	if err != nil {
		return nil, err
	}

	go messageHandler(cfg, messages, c.done)
	return c, nil
}

func (c consumer) Close() error {
	c.conn.Close()
	return <-c.done
}

func messageHandler(cfg *SiteConfiguration, messages <-chan amqp.Delivery, done chan error) {

	for msg := range messages {

		// store unmarshaled code in a temporary variable to prevent config to be corrupt if error occures
		unmarshaled := &SiteConfiguration{}
		err := json.Unmarshal(msg.Body, unmarshaled)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		*cfg = *unmarshaled

		msg.Ack(false)
	}
	done <- nil
}