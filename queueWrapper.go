package rabbitHelper

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

const (
	amqpDefaultHeartbeat = 10 * time.Second
	amqpDefaultLocale    = "en_US"
)

type QueueWrapper struct {
	Host  string
	User  string
	Pswd  string
	Vhost string
	Port  string

	queueName       string
	connection      *amqp.Connection
	channel         *amqp.Channel
	MessagesChannel <-chan amqp.Delivery
	queueDef        queueDefinitions
	exchangeDef     exchangeDefinitions
	bindDef         bindDefinitions
	publishDef      publishDefinitions
	consumeDef      consumeDefinitions
}

type queueDefinitions struct {
	shouldCreate bool
	durable      bool
	autoDelete   bool
	exclusive    bool
	noWait       bool
	args         amqp.Table
}

type exchangeDefinitions struct {
	shouldSet  bool
	name       string
	kind       string
	durable    bool
	autoDelete bool
	internal   bool
	noWait     bool
	args       amqp.Table
}

type bindDefinitions struct {
	shouldSet bool
	queueName string
	key       string
	exchange  string
	noWait    bool
}

type publishDefinitions struct {
	exchange    string
	routingKey  string
	mandatory   bool
	immediate   bool
	contentType string
}

type consumeDefinitions struct {
	consumer  string
	autoAck   bool
	exclusive bool
	noLocal   bool
	noWait    bool
}

var ErrQueueNotFound = errors.New("queue not found or unexpected definitions")
var ErrPublish = errors.New("failed to publish on channel")

func NewQueueWrapper(queueName string, opts ...Option) *QueueWrapper {
	res := &QueueWrapper{
		Host:  "localhost",
		User:  "guest",
		Pswd:  "guest",
		Vhost: "/",
		Port:  "5672",

		queueName: queueName,
		queueDef: queueDefinitions{
			shouldCreate: true,
			durable:      true,
			autoDelete:   false,
			exclusive:    false,
			noWait:       false,
		},
		exchangeDef: exchangeDefinitions{
			shouldSet:  false,
			name:       "",
			kind:       "",
			durable:    true,
			autoDelete: false,
			internal:   false,
			noWait:     false,
			args:       nil,
		},
		bindDef: bindDefinitions{
			shouldSet: false,
			queueName: "",
			key:       "",
			exchange:  "",
			noWait:    false,
		},
		publishDef: publishDefinitions{
			exchange:    "",
			routingKey:  queueName,
			mandatory:   false,
			immediate:   false,
			contentType: "application/json",
		},
		consumeDef: consumeDefinitions{
			consumer:  "",
			autoAck:   false,
			exclusive: false,
			noLocal:   false,
			noWait:    false,
		},
	}

	// apply config options
	for _, opt := range opts {
		opt(res)
	}

	return res
}

func (q *QueueWrapper) Close() {
	if err := q.channel.Close(); err != nil {
		log.Warn("failed to close rabbit channel ["+q.queueName+"]", err)
	}
	if err := q.connection.Close(); err != nil {
		log.Warn("failed to close rabbit connection ["+q.queueName+"]", err)
	}

}

func (q *QueueWrapper) Open() error {
	if err := q.createConnection(); err != nil {
		return err
	}
	if err := q.createChannel(); err != nil {
		return err
	}

	if q.exchangeDef.shouldSet {
		if err := q.makeSureExchangeExists(); err != nil {
			return err
		}
	}

	if qErr := q.makeSureQueueExists(); qErr != nil {
		log.Error(fmt.Sprintf("Rabbit queue not found or configured not as expected [%s]", q.queueName), qErr)
		return ErrQueueNotFound
	}

	if q.bindDef.shouldSet {
		if err := q.makeSureBindExists(); err != nil {
			return err
		}
	}

	return nil
}

func (q *QueueWrapper) makeSureExchangeExists() error {
	err := q.channel.ExchangeDeclare(
		q.exchangeDef.name,       // queueName
		q.exchangeDef.kind,       // type
		q.exchangeDef.durable,    // durable
		q.exchangeDef.autoDelete, // auto-deleted
		q.exchangeDef.internal,   // internal
		q.exchangeDef.noWait,     // no-wait
		q.exchangeDef.args,       // arguments
	)
	if err != nil {
		log.Error("Failed to declare exchange", err)
		return err
	}
	return nil
}

func (q *QueueWrapper) makeSureBindExists() error {
	err := q.channel.QueueBind(
		q.bindDef.queueName, // queue queueName
		q.bindDef.key,       // routing key
		q.bindDef.exchange,  // exchange
		q.bindDef.noWait,
		nil)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to bind %s to %s", q.bindDef.queueName, q.bindDef.exchange), err)
		return err
	}
	return nil
}

func (q *QueueWrapper) createConnection() error {
	conn, err := GetConnection(q.Host, q.User, q.Pswd, q.Vhost, q.Port)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}
	q.connection = conn
	return nil
}

func (q *QueueWrapper) createChannel() error {
	ch, err := q.connection.Channel()
	if err != nil {
		log.Error("Failed to open a channel: ", err)
		return err
	}
	q.channel = ch
	return nil
}

func (q *QueueWrapper) makeSureQueueExists() error {
	if !q.queueDef.shouldCreate {
		return nil
	}
	_, err := q.channel.QueueDeclare(
		q.queueName,           // queueName
		q.queueDef.durable,    // durable
		q.queueDef.autoDelete, // delete when unused
		q.queueDef.exclusive,  // exclusive
		q.queueDef.noWait,     // no-wait
		q.queueDef.args,       // arguments
	)
	if err != nil {
		log.Error("Failed to declare a queue: ", err.Error())
		return ErrQueueNotFound
	}
	err = q.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	// see: https://www.rabbitmq.com/tutorials/tutorial-two-go.html (Fair dispatch)
	if err != nil {
		log.Error("Failed to config queue prefetch to 1: ", err.Error())
		return err
	}

	return nil
}

func (q *QueueWrapper) StartConsuming() error {
	if qErr := q.makeSureQueueExists(); qErr != nil {
		log.Fatal(fmt.Sprintf("Rabbit queue not found or configured not as expected [%s]", q.queueName), qErr)
	}

	var err error
	q.MessagesChannel, err = q.channel.Consume(
		q.queueName,            // queue
		q.consumeDef.consumer,  // consumer
		q.consumeDef.autoAck,   // auto-ack
		q.consumeDef.exclusive, // exclusive
		q.consumeDef.noLocal,   // no-local
		q.consumeDef.noWait,    // no-wait
		nil,                    // args
	)
	return err
}

func (q *QueueWrapper) Publish(msg []byte) error {
	if qErr := q.makeSureQueueExists(); qErr != nil {
		log.Error(fmt.Sprintf("Rabbit queue not found or configured not as expected [%s]", q.queueName), qErr)
		return ErrQueueNotFound
	}

	err := q.channel.Publish(
		q.publishDef.exchange,   // exchange
		q.publishDef.routingKey, // routing key
		q.publishDef.mandatory,  // mandatory
		q.publishDef.immediate,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  q.publishDef.contentType,
			Body:         msg,
		})
	if err != nil {
		return ErrPublish
	}
	return nil
}

func GetConnection(host, user, pswd, vhost, port string) (*amqp.Connection, error) {
	config := amqp.Config{
		Vhost:     vhost,
		Heartbeat: amqpDefaultHeartbeat,
		Locale:    amqpDefaultLocale,
	}

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pswd, host, port)

	uri, err := amqp.ParseURI(dsn)
	if err != nil {
		return nil, err
	}

	if config.Vhost == "" {
		config.Vhost = uri.Vhost
	}

	return amqp.DialConfig(dsn, config)
}
