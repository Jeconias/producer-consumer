package manager

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Jeconias/producer-consumer/texterrors"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//NewWS new struct of ManagerAMQPWS
func NewWS() *ManagerAMQPWS {
	return &ManagerAMQPWS{
		conns:    make(map[string]*websocket.Conn, 0),
		amqpConn: StartDefault(),
	}
}

//AppendConn will add new socket in map
func (ws *ManagerAMQPWS) AppendConn(conn *websocket.Conn) {
	addr := conn.RemoteAddr().String()
	if ok := ws.HasConn(addr); !ok {
		ws.conns[addr] = conn
	}
}

//Conns will return a copy of all socket
func (ws *ManagerAMQPWS) Conns() map[string]*websocket.Conn {
	return ws.conns
}

func (ws *ManagerAMQPWS) AmqpConn() *Conn {
	return ws.amqpConn
}

//HasConn will verify if exist a socket in map with equal the key of parameter
func (ws *ManagerAMQPWS) HasConn(key string) bool {
	ok := false
	if _, ok = ws.conns[key]; ok {
		return ok
	}
	return ok
}

//DataMessage will get data massge of socket and return a struct of TransferData
func (ws *ManagerAMQPWS) DataMessage(key string) (*TransferData, error) {

	if ok := ws.HasConn(key); !ok {
		//TODO add pattern
		return nil, errors.New("No socket")
	}

	msgType, msg, err := ws.conns[key].ReadMessage()
	if err != nil {
		switch cErr := err.(type) {
		case *websocket.CloseError:
			if cErr.Code == 1001 {
				delete(ws.conns, key)
				return nil, errors.New(cErr.Error())
			}
		}
		if len(msg) == 0 {
			return nil, errors.New(texterrors.NoMessageToSend)
		}
		return nil, err
	}

	data := &TransferData{}
	if err := json.Unmarshal(msg, data); err != nil {
		return nil, err
	}
	data.MessageType = msgType
	data.From = ws.conns[key]
	return data, nil
}

func (ws *ManagerAMQPWS) HandleNewConnection(new *websocket.Conn) {
	addr := new.RemoteAddr().String()
	if ok := ws.HasConn(addr); !ok {
		fmt.Printf("New connection: %s\n", addr)
	}
	ws.AppendConn(new)
}

//Send will send the data for all sockets connected or to one
func (ws *ManagerAMQPWS) SendToQueue(from *websocket.Conn) error {

	tData, err := ws.DataMessage(from.RemoteAddr().String())
	if err != nil {
		return err
	}

	err = ws.ProducerQueue("pac", []byte(tData.Data))
	if err != nil {
		//TODO Add pattern to solve the panic
		return err
	}
	return nil
}

func (ws *ManagerAMQPWS) SendToAllWSClient(from *websocket.Conn) error {

	sSocket := from.RemoteAddr().String()
	tData, err := ws.DataMessage(sSocket)
	if err != nil {
		return err
	}

	for _, s := range ws.conns {
		if tData.To == "all" && s.RemoteAddr().String() != from.RemoteAddr().String() {
			err := s.WriteMessage(tData.MessageType, []byte(tData.Data))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ws *ManagerAMQPWS) ProducerQueue(queueName string, toSend []byte) error {
	if _, ok := ws.amqpConn.queues[queueName]; !ok {
		return errors.New(texterrors.QueueNotFound)
	}

	return ws.amqpConn.channel.Publish("", ws.amqpConn.queues[queueName].Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        toSend,
	})
}

func (ws *ManagerAMQPWS) ConsumeQueue(queueName string) error {

	loop := make(chan bool)

	msgs, err := ws.amqpConn.channel.Consume(ws.amqpConn.queues[queueName].Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			fmt.Printf("From consume: %s\n", d.Body)
		}
	}()

	<-loop

	return nil
}

func (ws *ManagerAMQPWS) Close() error {
	return ws.amqpConn.conn.Close()
}
