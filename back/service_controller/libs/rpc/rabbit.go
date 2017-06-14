package rpc

import (
	"github.com/streadway/amqp"
	"github.com/johnnovikov/hackforces/back/service_controller/libs/helpers"
	"fmt"
)

type RabbitMqRpc struct {
	Handler DataHandler
	Connection *amqp.Connection
	Channel *amqp.Channel
	Queue *amqp.Queue
	WorkersNum int
}

func (mq *RabbitMqRpc) Build(queueName string,workerNums int) {
	ch, err := mq.Connection.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	mq.Channel = ch
	q, err := ch.QueueDeclare(
		queueName, // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	mq.Queue = &q
	helpers.FailOnError(err, "Failed to declare a queue")
	err = ch.Qos(
				workerNums,     // prefetch count
                0,     // prefetch size
                false, // global
        )
	helpers.FailOnError(err, "Failed to set QoS")
	mq.WorkersNum = workerNums
	//defer ch.Close()

}

func (mq *RabbitMqRpc) Close() {
	fmt.Print("LOL")
	mq.Channel.Close()
}


func (mq *RabbitMqRpc) Handle() {
	msgs, err := mq.Channel.Consume(
                mq.Queue.Name, // queue
                "",     // consumer
                false,  // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
	)
	helpers.FailOnError(err,"Failed to register a consumer")
	for i:=0; i < mq.WorkersNum; i++{
	  go func() {
		for d := range msgs {
			request := string(d.Body)
			response := mq.Handler.HandleRequest(request)

			err = mq.Channel.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
							ContentType:   "text/plain",
							CorrelationId: d.CorrelationId,
							Body:          []byte(response),
					})
			helpers.FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	  }()
	}
	forever := make(chan bool)
	<-forever
}
