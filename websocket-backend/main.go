package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type TaskRequest struct {
	TaskID string `json:"task_id"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for simplicity
	},
}

func getTask(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var taskRequest TaskRequest
		if err := json.Unmarshal(message, &taskRequest); err != nil {
			log.Println("unmarshal:", err)
			break
		}
		
		// this string is printed as a response
		response := "Received task_id: " + taskRequest.TaskID
		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/get_task", getTask)
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
