package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Room string
}

type Hub struct {
	mu    sync.RWMutex
	rooms map[string]map[*Client]struct{}
}

func NewHub() *Hub { return &Hub{rooms: map[string]map[*Client]struct{}{}} }

func (h *Hub) Join(c *Client) {
	h.mu.Lock()
	if h.rooms[c.Room] == nil {
		h.rooms[c.Room] = map[*Client]struct{}{}
	}
	h.rooms[c.Room][c] = struct{}{}
	h.mu.Unlock()
}

func (h *Hub) Leave(c *Client) {
	h.mu.Lock()
	if m, ok := h.rooms[c.Room]; ok {
		delete(m, c)
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(room string, payload []byte) {
	h.mu.RLock()
	for c := range h.rooms[room] {
		select {
		case c.Send <- payload:
		default:
		}
	}
	h.mu.RUnlock()
}
