package rabbitHelper

import "github.com/streadway/amqp"

type Option func(*QueueWrapper)

func Host(host string) Option {
	return func(q *QueueWrapper) {
		q.Host = host
	}
}
func User(user string) Option {
	return func(q *QueueWrapper) {
		q.User = user
	}
}
func Pswd(pswd string) Option {
	return func(q *QueueWrapper) {
		q.Pswd = pswd
	}
}
func Vhost(vhost string) Option {
	return func(q *QueueWrapper) {
		q.Vhost = vhost
	}
}
func Port(port string) Option {
	return func(q *QueueWrapper) {
		q.Port = port
	}
}

func CreateQueue(create bool) Option {
	return func(q *QueueWrapper) {
		q.queueDef.shouldCreate = create
	}
}

func Queue(durable, autoDelete, exclusive, noWait bool, args amqp.Table) Option {
	return func(q *QueueWrapper) {
		q.queueDef.durable = durable
		q.queueDef.autoDelete = autoDelete
		q.queueDef.exclusive = exclusive
		q.queueDef.noWait = noWait
		q.queueDef.args = args
	}
}

func Publish(exchange, routingKey string, mandatory, immediate bool, contentType string) Option {
	return func(q *QueueWrapper) {
		q.publishDef.exchange = exchange
		q.publishDef.routingKey = routingKey
		q.publishDef.mandatory = mandatory
		q.publishDef.immediate = immediate
		q.publishDef.contentType = contentType
	}
}

func Consume(consumer string, autoAck, exclusive, noLocal, noWait bool) Option {
	return func(q *QueueWrapper) {
		q.consumeDef.consumer = consumer
		q.consumeDef.autoAck = autoAck
		q.consumeDef.exclusive = exclusive
		q.consumeDef.noLocal = noLocal
		q.consumeDef.noWait = noWait
	}
}

func Exchange(name string, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) Option {
	return func(q *QueueWrapper) {
		q.exchangeDef.shouldSet = true
		q.exchangeDef.name = name
		q.exchangeDef.kind = kind
		q.exchangeDef.durable = durable
		q.exchangeDef.autoDelete = autoDelete
		q.exchangeDef.internal = internal
		q.exchangeDef.noWait = noWait
		q.exchangeDef.args = args
	}
}

func Bind(queueName, routingKey, exchange string) Option {
	return func(q *QueueWrapper) {
		q.bindDef.shouldSet = true
		q.bindDef.queueName = queueName
		q.bindDef.key = routingKey
		q.bindDef.exchange = exchange
		q.bindDef.noWait = false
	}
}

func DeadLetter(exchange, queue, routing string) Option {
	return func(q *QueueWrapper) {
		q.deadLetter.shouldSet = true
		q.deadLetter.exchange = exchange
		q.deadLetter.routing = routing
		q.deadLetter.queue = queue
	}
}
