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
var broadcast = make(chan node.Message)
var upgrader = websocket.Upgrader{}

type Message struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
	Message		string `json:"message"`
}



func main() {
	log.Printf("my peer is " + os.Getenv("PEERS"))


	var node node.Node
	node = app.NewClient("8001", os.Getenv("PEERS"), broadcast)
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
		err := ws.ReadJSON(&msg)
		if err != nil{
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		log.Printf("no problem\n")
		broadcast <- msg
	}
}

func handleMessages(){
	var msg node.Message
	var jsonMessage Message
	for{
		msg = <- broadcast
		jsonMessage.Message = msg.Content
		log.Printf("UI received message $s\n", jsonMessage.Message)
		for client := range clients{
			log.Printf("sending message...\n")
			err := client.WriteJSON(jsonMessage)
			if err != nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
