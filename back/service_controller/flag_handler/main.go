package main

//import "github.com/johnnovikov/hackforces/back/service_controller/libs"
import (
	//"github.com/streadway/amqp"
	"github.com/johnnovikov/hackforces/back/service_controller/libs/rpc"
	//"net/http"
	"github.com/streadway/amqp"
	"github.com/johnnovikov/hackforces/back/service_controller/libs/helpers"
	"sync"
)


//func handle(data string)

type TeamData struct {
	id int
	points int
	mu sync.Mutex
}

func NewTeamData(id int,points int) *TeamData{
	return &TeamData{id:id,points:points,mu:sync.Mutex{}}
}




func main()  {
	conn , err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	helpers.FailOnError(err,"Cant connect to rabbit")
	defer conn.Close()
	a_handler:= new(rpc.AckHandler)
	a_handler.Init()

	mq := new(rpc.RabbitMqRpc)
	defer mq.Close()
	mq.Connection = conn
	mq.Build("flags_rpc",1)
	mq.Handler = a_handler
	mq.Handle()
}


