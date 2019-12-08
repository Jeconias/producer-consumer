package main

import (
	"fmt"
	"net/http"

	"github.com/Jeconias/producer-consumer/manager"

	"github.com/gorilla/websocket"
)

var (
	upgrader websocket.Upgrader
	sockets  *manager.ManagerAMQPWS
)

func handleWS(writer http.ResponseWriter, request *http.Request) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		//TODO Add pattern to solve the panic
		panic(err)
	}

	sockets.HandleNewConnection(socket)
	for {
		if err := sockets.SendToQueue(socket); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}

func main() {

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	sockets = manager.NewWS()
	//TODO Remove Consume of main
	go sockets.ConsumeQueue("pac")

	http.HandleFunc("/", handleWS)
	http.ListenAndServe(":8000", nil)

}
