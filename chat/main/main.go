package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"LOG735-PG/src/node"
	"LOG735-PG/src/app"
	"time"
	"os"
)

var clients = make(map[*websocket.Conn]bool)
var uiChannel = make(chan node.Message)
var nodeChannel = make(chan node.Message)
var upgrader = websocket.Upgrader{}

type Message struct {
	Peer		string `json:"peer"`
	Message		string `json:"message"`
}

func main() {



	log.Printf("my peer is " + os.Getenv("PEERS"))

	var node node.Node
	node = app.NewClient("8001", os.Getenv("PEERS"), uiChannel, nodeChannel)
	err := node.SetupRPC("8001")
	if err != nil {
		log.Fatal("RPC setup error:", err)
	}
	err = node.Peer()
	if err != nil {
		log.Fatal("Peering error:", err)
	}

	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	log.Println("http server started on :8000")

	go func() {
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	for {
		time.Sleep(time.Hour)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request){

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil{
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true



	for{
		var msg node.Message
		var appMsg Message
		err := ws.ReadJSON(&appMsg)
		msg.Content = appMsg.Message
		msg.Peer = appMsg.Peer
		if err != nil{
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		uiChannel <- msg
		nodeChannel <- msg
	}
}

func handleMessages(){
	var msg node.Message
	var jsonMessage Message
	for{
		msg = <- uiChannel
		jsonMessage.Message = msg.Content
		jsonMessage.Peer = msg.Peer
		for client := range clients{
			err := client.WriteJSON(jsonMessage)
			if err != nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
