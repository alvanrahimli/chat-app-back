package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	messengerWs "tree-messenger/pkg/messenger-ws"
)

type CommandType string
const (
	Registration CommandType = "register"
	BecomeOnline CommandType = "become_alive"
)

type ClientCommand struct {
	Type	CommandType `json:"type"`
	Content string `json:"content"`
}

type ClientResponse struct {
	Status	string `json:"status"`
	Content string `json:"content"`
}

func serveWs(pool *messengerWs.Pool, w http.ResponseWriter, r *http.Request) {
	wsConn, upgradeErr := messengerWs.Upgrade(w, r)
	if upgradeErr != nil {
		log.Printf("Could not upgrade connection\n\t %s", upgradeErr.Error())
		return
	}

	for {
		var clientCommand ClientCommand

		_, p, msgErr := wsConn.ReadMessage()
		if msgErr != nil {
			log.Printf("Error occured whire reading message. Error: %s", msgErr.Error())
			continue
		}

		jsonErr := json.Unmarshal(p, &clientCommand)
		if jsonErr != nil {
			log.Printf("Error: %s", jsonErr.Error())
			continue
		}

		if clientCommand.Type == Registration {
			client := messengerWs.Client{
				ID:			uuid.New().String(),
				Name:		clientCommand.Content,
				Conn:		wsConn,
				LastPing: 	time.Now().UTC(),
			}

			responseObj := ClientResponse {
				Status:  "ok",
				Content: client.ID,
			}

			responseJson, marshallErr := json.Marshal(responseObj)
			if marshallErr != nil {
				log.Printf("Error: %s", marshallErr.Error())
				continue
			}

			responseErr := wsConn.WriteMessage(websocket.TextMessage, responseJson)
			if responseErr != nil {
				log.Printf("RESPONSE ERROR: %s", responseErr.Error())
				continue
			}

			pool.Register <- &client
			client.Read()
		} else if clientCommand.Type == BecomeOnline {
			pool.MakeOnline <- clientCommand.Content
		}

	}
}

func setupRoutes() {
	pool := messengerWs.InitializePool()
	go pool.Start()

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(pool, writer, request)
	})
}

func main() {
	fmt.Println("Messenger v0.1")

	setupRoutes()

	log.Println("Application started, listening port :8000")
	httpError := http.ListenAndServe(":8000", nil)
	if httpError != nil {
		log.Printf("Error: %s", httpError.Error())
	}
}
