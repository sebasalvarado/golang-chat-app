package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

type room struct {
	// forward is a channel that holds incomming messages
	// that should be forwarded to other clients
	forward chan []byte
	// join is a channel for clients wishing to join the room
	join chan *client
	// leave is a channel for cleints wishing to leave the room
	leave chan *client

	// clients holds all current clients in this room
	// We need join and leave channels to safely add/remove clients from room
	// o/w concurrent threads might want to access the map
	clients map[*client]bool
}

// Create a new instance of type room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
func (r *room) run() {
	// select statements: Useful for synchronizing or modifying shared memory or take different actions depending
	// on the various activities within our channels
	for {
		// This inifite for loop will keep watching all three channels from type room: join, leave and forward.
		// If there is any change in any of the channels the correct select case will be run. ONly ONE block at a time
		// that way it can synchronize the execution
		select {
		case client := <-r.join:
			// joining
			r.clients[client] = true
		case client := <-r.leave:
			// leaving the room
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// forward messages to all connected clients
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Upgrades the HTTP Connection to get a WebSocket
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	// Initialize the client
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()

	// go: Initialize a new thread and execute that code in that new thread
	go client.write()
	client.read()
}
