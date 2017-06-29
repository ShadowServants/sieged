package rpc

import (
	"github.com/streadway/amqp"
	"github.com/jnovikov/hackforces/back/libs/helpers"
	"fmt"
	"sync"

)

type RabbitMqRpc struct {
	Handler         DataHandler
	Connection      *amqp.Connection
	Channel         *amqp.Channel
	Queue           *amqp.Queue
	WorkersNum      int
	messages        <- chan amqp.Delivery
	parallelMethods chan bool
	wg              sync.WaitGroup
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
	mq.parallelMethods = make(chan bool, mq.WorkersNum)
	//defer ch.Close()
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
	mq.messages = msgs
}

func (mq *RabbitMqRpc) Close() {
	fmt.Print("LOL")
	mq.Channel.Cancel("",false)
	mq.wg.Wait()
	mq.Channel.Close()

}



func (mq *RabbitMqRpc) Handle() {
	for d := range mq.messages {
		mq.parallelMethods <- true //Это просто семафор. Не трогай его и он не тронет тебя
		mq.wg.Add(1) //Для safe-way завершения процесса
		go mq.handleRequest(d)
	}
	mq.wg.Wait()
}

func (mq *RabbitMqRpc) handleRequest(d amqp.Delivery) {
	defer func() {
		<-mq.parallelMethods
		mq.wg.Done()
	}()
	request := string(d.Body)
	response := mq.Handler.HandleRequest(request)
	err := mq.Channel.Publish(
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
