package main

import (
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {
	manager := socket.NewManager()
	http.HandleFunc("/ws", manager.ServeWS)
}
