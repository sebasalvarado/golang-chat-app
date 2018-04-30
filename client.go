package main

import (
	"github.com/gorilla/websocket"
)

// Client represents a single chatting user

type client struct {
	// socket is the web socket that this client has
	socket *websocket.Conn
	// Send is channel in which messages are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

// Channels: Are an in-memory thread safe queue. Where senders pass data and receivers read data
// in a non-blocking, thread-safe fashion

// Methods that allows us to read from the socket and forwards it to the forward channel on the room type
func (c *client) read() {
	defer c.socket.Close() // defer runs at the end of the function
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
