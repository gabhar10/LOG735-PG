package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"time"
	"net/rpc"
	//"os"
	//brpc "LOG735-PG/src/rpc"
	"os"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var connections []PeerConnection

var upgrader = websocket.Upgrader{}

type Message struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
	Message		string `json:"message"`
}

type PeerConnection struct {
	ID			string
	conn		*rpc.Client
}

func main() {
	log.Printf("my peer is " + os.Getenv("PEERS"))

	/*client, err := brpc.ConnectTo(os.Getenv("PEERS"))
	if err != nil {
		log.Printf("Woops, error when connecting to client node\n")
	}
	args := &brpc.ConnectionRPC{"8000"}
	var reply brpc.BlocksRPC
	err = client.Call("NodeRPC.Peer", args, &reply)
	var newConnection = new(PeerConnection)
	newConnection.ID = os.Getenv("PEERS")
	newConnection.conn = client
	connections = append(connections, *newConnection)*/


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
		var msg Message

		err := ws.ReadJSON(&msg)
		if err != nil{
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages(){
	for{
		msg := <- broadcast
		for client := range clients{
			err := client.WriteJSON(msg)
			if err != nil{
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
