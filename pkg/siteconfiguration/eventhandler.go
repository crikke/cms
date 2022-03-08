package siteconfiguration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/streadway/amqp"
)

type ConfigurationEventHandler interface {
	io.Closer
	Watch(cfg *SiteConfiguration) error
	Publish(cfg SiteConfiguration) error
}

const cfgExchange = "cms.siteconfiguration"

type configurationQueue struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       string
	tag         string
	done        chan error
	initialized bool
}

// Initializes a temporary queue that subscribes to configuration changes
func NewConfigurationEventHandler(uri string) (ConfigurationEventHandler, error) {
	c := &configurationQueue{
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
	c.queue = q.Name

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

	c.initialized = true

	return c, nil
}

func (c *configurationQueue) Close() error {
	c.conn.Close()
	c.initialized = false
	return <-c.done
}

func (c configurationQueue) Watch(cfg *SiteConfiguration) error {

	if !c.initialized {
		return errors.New("cannot watch before channel is initialized")
	}

	messages, err := c.channel.Consume(c.queue, "", false, true, false, false, nil)

	if err != nil {
		return err
	}

	go messageHandler(cfg, messages, c.done)

	return nil
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

func (c configurationQueue) Publish(cfg SiteConfiguration) error {

	data, err := json.Marshal(&cfg)

	if err != nil {
		return err
	}

	return c.channel.Publish(
		cfgExchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
}
