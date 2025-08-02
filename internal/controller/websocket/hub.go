package websocket

import (
	"encoding/json"
	"sync"
)

type SignalingMessage struct {
	Type    string          `json:"type"`
	To      string          `json:"to"`
	From    string          `json:"from"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	clients    map[string]*Peer
	broadcast  chan SignalingMessage
	register   chan *Peer
	unregister chan *Peer
	mux        sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Peer),
		broadcast:  make(chan SignalingMessage),
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case register := <-h.register:
			h.Register(register)
		case unregister := <-h.unregister:
			h.Unregister(unregister)
		case message := <-h.broadcast:
			h.RouteMessage(message)
		}
	}
}

func (h *Hub) Register(peer *Peer) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.clients[peer.ID] = peer
}

func (h *Hub) Unregister(peer *Peer) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if _, ok := h.clients[peer.ID]; ok {
		delete(h.clients, peer.ID)
		close(peer.send)
	}
}

func (h *Hub) RouteMessage(message SignalingMessage) {
	h.mux.RLock()
	defer h.mux.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return
	}

	reciever := h.clients[message.To]
	if reciever != nil {
		reciever.send <- messageBytes
	}
}
