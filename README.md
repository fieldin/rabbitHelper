# Rabbit helper
This package wraps the official amqp package and provider some helpers for using it
Mainly it provides a wrapper for the queue abstraction.

I tried to provide reasonable defaults to all the settings, and applied the options pattern to change them.

## usage examples

### Simple consuming
```go
rbt := rabbitHelper.NewQueueWrapper("my_queue")
if err := rbt.Open(); err != nil {
    log.WithError(err).Fatal("failed opening incoming rabbit connection")
}
defer rbt.Close()

if err := rbt.StartConsuming(); err != nil {
    log.WithError(err).Fatal("failed consuming incoming messages")
}

for {
    select {
    case <-done:
        return
    case rm := <-rbt.MessagesChannel:
        rawBody := rm.Body
    .
    .
    .
```

### Consuming example, using many of the config options
```go
queueName := getQueueName()
exchangeName := viper.GetString("rabbit_exchange")
routingKey := viper.GetString("rabbit_routing_key")

user := viper.GetString("rabbit_user")
pswd := viper.GetString("rabbit_pswd")
host := viper.GetString("rabbit_host")
port := viper.GetString("rabbit_port")
vhost := viper.GetString("rabbit_vhost")

rbt := rabbitHelper.NewQueueWrapper(
    queueName,
    rabbitHelper.User(user), rabbitHelper.Pswd(pswd), rabbitHelper.Host(host), rabbitHelper.Vhost(vhost), rabbitHelper.Port(port),
    rabbitHelper.Exchange(exchangeName, "direct", true, false, false, false, nil),
    rabbitHelper.Bind(queueName, routingKey, exchangeName),
)

if err := rbt.Open(); err != nil {
    log.WithError(err).Fatal("failed opening incoming rabbit connection")
}
defer rbt.Close()

if err := rbt.StartConsuming(); err != nil {
    log.WithError(err).Fatal("failed consuming incoming messages")
}

for {
    select {
    case <-done:
        return
    case rm := <-rbt.MessagesChannel:
        rawBody := rm.Body

    .
    .
    .
```

### Publishing example
Since we didn't use the Exchange option, this will publish to the default exchange
```go
rbt := rabbitHelper.NewQueueWrapper(queueName)
if err := rbt.Open(); err != nil {
    log.WithError(err).Fatal("failed connection to rabbit")
}
defer rbt.Close()

publishError := rbt.Publish([]byte(job.String()))
if (publishError == rabbitHelper.ErrQueueNotFound) || (publishError == rabbitHelper.ErrPublish) {
    log.Fatal(publishError)
}
```