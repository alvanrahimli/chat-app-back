package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	messengerWs "tree-messenger/pkg/messenger-ws"

	"github.com/google/uuid"
)



func serveWs(pool *messengerWs.Pool, w http.ResponseWriter, r *http.Request) {
	wsConn, upgradeErr := messengerWs.Upgrade(w, r)
	if upgradeErr != nil {
		log.Printf("Could not upgrade connection\n\t %s", upgradeErr.Error())
		return
	}

	for {
		var clientCommand messengerWs.ClientCommand

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

		if clientCommand.Type == messengerWs.Registration {
			client := messengerWs.Client{
				ID:			uuid.New().String(),
				Name:		clientCommand.Content,
				Conn:		wsConn,
				LastPing: 	time.Now().UTC(),
				Pool:		pool,
			}

			responseObj := messengerWs.ClientResponse {
				Status:  "ok",
				Type: messengerWs.Registered,
				Content: client.ID,
			}

			sendErr := client.Send(responseObj)
			if sendErr != nil {
				log.Printf("Error: %s", sendErr.Error())
			}

			pool.Register <- &client
			client.Read()
		} else if clientCommand.Type == messengerWs.BecomeOnline {
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
	//httpError := http.ListenAndServeTLS(":8000", "/home/alvan/mkcert/localhost.local.pem", "/home/alvan/mkcert/localhost.local-key.pem", nil)
	httpError := http.ListenAndServe(":8000", nil)
	if httpError != nil {
		log.Printf("Error: %s", httpError.Error())
	}
}
