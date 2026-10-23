package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// SSEEvent represents a message sent to connected clients
type SSEEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SSEHub manages active client connections and broadcasts
type SSEHub struct {
	clients   map[chan string]bool
	broadcast chan string
	register  chan chan string
	unregister chan chan string
	mu        sync.RWMutex
}

var Hub *SSEHub

func init() {
	Hub = NewSSEHub()
	go Hub.Run()
}

// NewSSEHub creates a new SSEHub
func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:    make(map[chan string]bool),
		broadcast:  make(chan string, 100),
		register:   make(chan chan string),
		unregister: make(chan chan string),
	}
}

// Run listens for client registrations and broadcasts
func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("[SSE] Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
			}
			h.mu.Unlock()
			log.Printf("[SSE] Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client <- message:
				default:
					// Drop if buffer full to avoid blocking
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register registers a new client channel
func (h *SSEHub) Register(client chan string) {
	h.register <- client
}

// Unregister unregisters a client channel
func (h *SSEHub) Unregister(client chan string) {
	h.unregister <- client
}

// Broadcast sends an event to all connected clients
func (h *SSEHub) Broadcast(eventType string, payload interface{}) {
	event := SSEEvent{
		Type:    eventType,
		Payload: payload,
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[SSE] Error marshaling event: %v", err)
		return
	}
	msg := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(data))
	h.broadcast <- msg
}
