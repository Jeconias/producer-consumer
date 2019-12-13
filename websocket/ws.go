package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Jeconias/producer-consumer/definitions"
	"github.com/Jeconias/producer-consumer/texterrors"
	"github.com/gorilla/websocket"
)

//Default settings
var (
	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	//WsManager will save all connections of type WebSocket
	WsManager *Ws
)

//Ws is a struct with connections and func to handle connections of type WebSocket
type Ws struct {
	conns map[string]*websocket.Conn
}

//StartWs will starting the global variable WsManager
func StartWs() *Ws {
	WsManager = &Ws{
		conns: make(map[string]*websocket.Conn, 0),
	}
	return WsManager
}

//AppendConn will add new socket in map
func (ws *Ws) AppendConn(conn *websocket.Conn) {
	addr := conn.RemoteAddr().String()
	if ok := ws.HasConn(addr); !ok {
		ws.conns[addr] = conn
	}
}

//Conns will return all sockets
func (ws *Ws) Conns() map[string]*websocket.Conn {
	return ws.conns
}

//HasConn will verify if exist a socket in map with equal the key of parameter
func (ws *Ws) HasConn(key string) bool {
	ok := false
	if _, ok = ws.conns[key]; ok {
		return ok
	}
	return ok
}

//RemoveSocket will delete connection of map conns
func (ws *Ws) RemoveSocket(socket *websocket.Conn) {
	delete(ws.conns, socket.RemoteAddr().String())
}

//HandleNewConnection will add new connection in map if not exist and show informations
func (ws *Ws) HandleNewConnection(new *websocket.Conn) {
	addr := new.RemoteAddr().String()
	if ok := ws.HasConn(addr); !ok {
		fmt.Printf("New connection: %s\n", addr)
	}
	ws.AppendConn(new)
}

//DataMessage will get data massge of socket and return a struct of TransferData
func (ws *Ws) DataMessage(socket *websocket.Conn) (int, []byte, error) {
	msgType, msg, err := socket.ReadMessage()
	if err != nil {
		switch cErr := err.(type) {
		case *websocket.CloseError:
			if cErr.Code == 1001 {
				ws.RemoveSocket(socket)
				return 0, []byte(""), errors.New(cErr.Error())
			}
		}
		if len(msg) == 0 {
			return 0, []byte(""), errors.New(texterrors.NoMessageToSend)
		}
		return 0, []byte(""), err
	}
	return msgType, msg, err
}

func (ws *Ws) TDataWithFrom(from *websocket.Conn) (*definitions.TransferData, error) {
	_, msg, err := ws.DataMessage(from)
	if err != nil {
		return nil, err
	}

	data := &definitions.TransferData{}
	if err := json.Unmarshal(msg, data); err != nil {
		return nil, err
	}
	data.From = from.RemoteAddr().String()
	return data, nil
}

//SendToAll send data for all connections in WsManager
func (ws *Ws) SendToAll(from string, data string) error {
	for _, s := range ws.conns {
		if s.RemoteAddr().String() == from {
			continue
		}
		err := s.WriteMessage(1, []byte(data))
		if err != nil {
			return err
		}
	}
	return nil
}
