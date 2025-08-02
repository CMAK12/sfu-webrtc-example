package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Peer struct {
	ID   string
	Conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

func NewPeer(id string, conn *websocket.Conn, hub *Hub) *Peer {
	return &Peer{
		ID:   id,
		Conn: conn,
		send: make(chan []byte, 256),
		hub:  hub,
	}
}

func (p *Peer) ReadMessages() {
	defer func() {
		p.hub.Unregister(p)
		p.Conn.Close()
	}()

	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			return
		}

		var signal SignalingMessage
		if err := json.Unmarshal(message, &signal); err != nil {
			continue
		}

		signal.From = p.ID
		p.hub.broadcast <- signal
	}
}

func (p *Peer) WriteMessages() {
	defer p.Conn.Close()

	for message := range p.send {
		if err := p.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
