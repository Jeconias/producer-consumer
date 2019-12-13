package pac

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Jeconias/producer-consumer/definitions"
	"github.com/Jeconias/producer-consumer/texterrors"
	"github.com/Jeconias/producer-consumer/websocket"
	"github.com/streadway/amqp"
)

//Conn struct for manager conn, channel and queues of ReabbitMQ
type Conn struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue
}

//ConfigQueue is struct for add config in queue
type ConfigQueue struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

var (
	//ManagerAMQP is varaibel global for RabbitMQ
	ManagerAMQP *Conn
	err         error
)

//StartManagerAMQP is method for start the variable global ManagerAMQP
func StartManagerAMQP() *Conn {

	ManagerAMQP, err = NewConn("amqp://rabbitmq:rabbitmq@localhost:3000")

	if err != nil {
		//TODO Add pattern to solve the panic
		panic(err)
	}

	err = ManagerAMQP.CreateQueue("pac", ConfigQueue{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	})

	if err != nil {
		//TODO Add pattern to solve the panic
		panic(err)
	}

	//start consumer
	go func() {
		err = ManagerAMQP.ConsumeQueue("pac")
		if err != nil {
			//TODO Add pattern to solve the panic
			panic(err)
		}
	}()

	return ManagerAMQP
}

// NewConn create the new connection of amqp
func NewConn(addressAmpq string) (*Conn, error) {
	amqpConn, err := amqp.Dial(addressAmpq)
	if err != nil {
		return nil, errors.New(texterrors.FailedConnectionWithRabbitMQ)
	}

	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, errors.New(texterrors.FailedToOpenChannel)
	}

	return &Conn{
		conn:    amqpConn,
		channel: ch,
		queues:  make(map[string]amqp.Queue, 0),
	}, nil
}

func (c *Conn) TData(msg []byte) (*definitions.TransferData, error) {

	Tdata := &definitions.TransferData{}
	if err := json.Unmarshal(msg, Tdata); err != nil {
		return nil, err
	}
	return Tdata, nil
}

//ProducerQueue will send data for queue
func (c *Conn) ProducerQueue(queueName string, toSend []byte) error {

	if _, ok := c.queues[queueName]; !ok {
		return errors.New(texterrors.QueueNotFound)
	}

	return c.channel.Publish("", c.queues[queueName].Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        toSend,
	})
}

//ConsumeQueue will receive all data of queue
func (c *Conn) ConsumeQueue(queueName string) error {

	loop := make(chan bool)

	msgs, err := c.channel.Consume(c.queues[queueName].Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			tData, err := c.TData(d.Body)
			if err != nil {
				fmt.Printf("Failed Unmarshal: %s\n", err)
			}
			websocket.WsManager.SendToAll(tData.From, tData.Data)
			fmt.Printf("From: %s | To: %s | data %s\n", tData.From, tData.To, tData.Data)
		}
	}()

	<-loop

	return nil
}

//CreateQueue create new Queue
func (c *Conn) CreateQueue(name string, config ConfigQueue) error {
	queue, err := c.channel.QueueDeclare(name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Args)
	if err != nil {
		return err
	}
	c.queues[name] = queue
	return nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
