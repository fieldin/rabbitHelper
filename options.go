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

func Queue(durable, autoDelete, exclusive, noWait bool, args amqp.Table) Option {
	return func(q *QueueWrapper) {
		q.queueDef = queueDefinitions{
			durable:    durable,
			autoDelete: autoDelete,
			exclusive:  exclusive,
			noWait:     noWait,
			args:       args,
		}
	}
}

func Publish(exchange, routingKey string, mandatory, immediate bool, contentType string) Option {
	return func(q *QueueWrapper) {
		q.publishDef = publishDefinitions{
			exchange:    exchange,
			routingKey:  routingKey,
			mandatory:   mandatory,
			immediate:   immediate,
			contentType: contentType,
		}
	}
}

func Consume(consumer string, autoAck, exclusive, noLocal, noWait bool) Option {
	return func(q *QueueWrapper) {
		q.consumeDef = consumeDefinitions{
			consumer:  consumer,
			autoAck:   autoAck,
			exclusive: exclusive,
			noLocal:   noLocal,
			noWait:    noWait,
		}
	}
}

func Exchange(name string, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) Option {
	return func(q *QueueWrapper) {
		q.exchangeDef = exchangeDefinitions{
			shouldSet:  true,
			name:       name,
			kind:       kind,
			durable:    durable,
			autoDelete: autoDelete,
			internal:   internal,
			noWait:     noWait,
			args:       args,
		}
	}
}

func Bind(queueName, routingKey, exchange string) Option {
	return func(q *QueueWrapper) {
		q.bindDef = bindDefinitions{
			shouldSet: true,
			queueName: queueName,
			key:       routingKey,
			exchange:  exchange,
			noWait:    false,
		}
	}
}
