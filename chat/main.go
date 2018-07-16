package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"LOG735-PG/src/node"
	brpc "LOG735-PG/src/rpc"
	"LOG735-PG/chat/ui"
	"time"
	"os"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}
var rpcHandler *brpc.NodeRPC

type Message struct {
	Email		string `json:"email"`
	Username	string `json:"username"`
	Message		string `json:"message"`
}



func main() {
	log.Printf("my peer is " + os.Getenv("PEERS"))


	var node node.Node
	node = ui.NewUiNode("8001", os.Getenv("PEERS"))
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
