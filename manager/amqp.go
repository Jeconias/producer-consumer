package manager

import (
	"errors"

	"github.com/Jeconias/producer-consumer/texterrors"

	"github.com/streadway/amqp"
)

func StartDefault() *Conn {
	conn, err := NewConn("amqp://rabbitmq:rabbitmq@localhost:3000")
	if err != nil {
		//TODO Add pattern to solve the panic
		panic(err)
	}

	err = conn.CreateQueue("pac", ConfigQueue{
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

	return conn
}

func (c *Conn) CreateQueue(name string, config ConfigQueue) error {
	queue, err := c.channel.QueueDeclare(name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Args)
	if err != nil {
		return err
	}
	c.queues[name] = queue
	return nil
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
