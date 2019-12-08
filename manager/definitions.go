package manager

import (
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

type Conn struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queues  map[string]amqp.Queue
}

type ConfigQueue struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ManagerAMQPWS struct {
	conns    map[string]*websocket.Conn
	amqpConn *Conn
}

type TransferData struct {
	Data        string `json:"data"`
	To          string `json:"to"`
	From        *websocket.Conn
	MessageType int
}
