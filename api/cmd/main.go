package main

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/control"
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {
	controller := control.NewController()

	http.HandleFunc("/ws", controller.SocketManager.ServeWS)
}
