package main

import (
	"encoding/json"
	"net/http"

	"github.com/Jeconias/producer-consumer/pac"
	"github.com/Jeconias/producer-consumer/websocket"
)

func handleWS(writer http.ResponseWriter, request *http.Request) {
	socket, err := websocket.Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		//TODO Add pattern to solve the panic
		panic(err)
	}
	websocket.WsManager.HandleNewConnection(socket)
	for {
		tData, err := websocket.WsManager.TDataWithFrom(socket)
		if err != nil {
			panic(err)
		}
		data, err := json.Marshal(tData)
		if err != nil {
			//TODO Add pattern to solve the panic
			panic(err)
		}
		err = pac.ManagerAMQP.ProducerQueue("pac", data)
		if err != nil {
			//TODO Add pattern to solve the panic
			panic(err)
		}
	}
}

func main() {

	websocket.StartWs()
	pac.StartManagerAMQP()

	http.HandleFunc("/", handleWS)
	http.ListenAndServe(":8000", nil)

}
