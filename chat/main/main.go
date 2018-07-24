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

const APP_PORT = "8001"
var clients = make(map[*websocket.Conn]bool)// Websocket Slice for exchange between server and web application
var uiChannel = make(chan node.Message) 	// Channel for incoming message from Client Node
var nodeChannel = make(chan node.Message) 	// Channel for outgoing message from web application
var upgrader = websocket.Upgrader{}			// Used to upgrade HTTP connection to websocket

type Message struct {
	Peer		string `json:"peer"`
	Message		string `json:"message"`
}

// Create client node and wait for connection from port 8000
func main() {
	var node node.Node
	node = app.NewClient(APP_PORT, os.Getenv("PEERS"), uiChannel, nodeChannel)
	err := node.SetupRPC(APP_PORT)
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

// Handler for new request to server and for messages coming from web application
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

// Handler for all incoming message from uiChannel (Message from node)
func handleMessages(){
	var msg node.Message
	var jsonMessage Message
	for{
		msg = <- uiChannel
		jsonMessage.Message = msg.Content
		jsonMessage.Peer = msg.Peer
		for client := range clients{
			// Send message to web application
			err := client.WriteJSON(jsonMessage)
			if err != nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
