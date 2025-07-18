package control

import (
	"encoding/json"
	"github.com/theAnuragMishra/mnnit-chess-club/api/internal/socket"
)

func RoomChange(c *Controller, event socket.Event, client *socket.Client) error {
	c.SocketManager.Lock()
	defer c.SocketManager.Unlock()
	var payload RoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	if client.Room == payload.RoomID {
		return nil
	}
	delete(c.SocketManager.Rooms[client.Room], client)
	client.Room = payload.RoomID
	if c.SocketManager.Rooms[payload.RoomID] == nil {
		c.SocketManager.Rooms[payload.RoomID] = make(map[*socket.Client]bool)
	}
	c.SocketManager.Rooms[payload.RoomID][client] = true
	// printing clients for debugging purposes
	// for _, x := range c.SocketManager.Rooms {
	// 	for y := range x {
	// 		fmt.Println(y)
	// 	}
	// }
	return nil
}

func LeaveRoom(c *Controller, event socket.Event, client *socket.Client) error {
	c.SocketManager.Lock()
	defer c.SocketManager.Unlock()
	var payload RoomPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return err
	}
	delete(c.SocketManager.Rooms[client.Room], client)
	if len(c.SocketManager.Rooms[client.Room]) == 0 {
		delete(c.SocketManager.Rooms, client.Room)
	}
	client.Room = ""
	// for _, x := range c.SocketManager.Rooms {
	// 	for y := range x {
	// 		fmt.Println(y)
	// 	}
	// }
	return nil
}
