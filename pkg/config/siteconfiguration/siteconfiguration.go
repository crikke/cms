package siteconfiguration

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"golang.org/x/text/language"
)

const cfgExchange = "cms.siteconfiguration"

// Configuration bound to site, such as root page & configured languages.
// Since this configuration is configured by users. It should not be stored as a ConfigMap.
// TODO: this can wait and have hardcoded defaults for now.
type Configuration struct {
	// Languages are configured by contentdelivery api. The elements are prioritized.
	Languages []language.Tag
	RootPage  uuid.UUID
}

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

// Initializes a temporary queue that subscribes to configuration changes
func NewSubscriber(uri string) error {
	c := &consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
		done:    make(chan error),
	}

	var err error
	c.conn, err = amqp.Dial(uri)

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
		cfgExchange,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		return err
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
		return err
	}

	err = c.channel.QueueBind(
		q.Name,
		"",
		cfgExchange,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(q.Name, "", false, true, false, false, nil)

	if err != nil {
		return err
	}

	go messageHandler(deliveries, c.done)
	return nil
}

func messageHandler(deliveries <-chan amqp.Delivery, done chan error) {

	for d := range deliveries {

		fmt.Printf("got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body)
	}
	done <- nil
}
