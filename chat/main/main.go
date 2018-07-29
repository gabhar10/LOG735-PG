package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"LOG735-PG/src/node"
	"LOG735-PG/src/client"
	"time"
	"os"
	"strings"
	"fmt"
)

const HttpPort = ":8000"

var clients = make(map[*websocket.Conn]bool)// Websocket Slice for exchange between server and web application
var uiChannel = make(chan node.Message) 	// Channel for incoming message from Client Node
var nodeChannel = make(chan node.Message) 	// Channel for outgoing message from web application
var upgrader = websocket.Upgrader{}			// Used to upgrade HTTP connection to websocket
var clientNode node.Node					// Server's client nod reference

type Message struct {
	Peer		string `json:"peer"`
	Message		string `json:"message"`
}

// Create client node and wait for connection from port 8000
func main() {
	peers := []node.Peer{}
	for _, s := range strings.Split(os.Getenv("PEERS"), " ") {
		p := node.Peer{
			Host: fmt.Sprintf("node-%s", s),
			Port: s}

		peers = append(peers, p)
	}

	clientNode = client.NewClient(os.Getenv("PORT"), peers, uiChannel, nodeChannel)
	err := clientNode.SetupRPC(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("RPC setup error:", err)
	}
	err = clientNode.Peer()
	if err != nil {
		log.Fatal("Peering error:", err)
	}

	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/disconnect", handleDisconnect)
	http.HandleFunc("/connect", handleConnect)
	http.HandleFunc("/getID", handleGetID)
	go handleMessages()


	go func() {
		err := http.ListenAndServe(HttpPort, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	clientNode.Start()

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

// Handler if user want's to Disconnect
func handleDisconnect(w http.ResponseWriter, r *http.Request){
	err := clientNode.Disconnect()
	if err != nil {
		http.Error(w, "Error during disconnection", 500)
		return
	}else{
		fmt.Fprint(w, " ")
	}
}

// Handler if user wants to connect to Anchor Miner
func handleConnect(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	field := r.FormValue("anchor")
	if field == os.Getenv("PORT") {
		http.Error(w, "Cannot connect to yourself", 500)
		return
	}
	clientNode.Connect(field)
}

func handleGetID(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, os.Getenv("PORT"))
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
