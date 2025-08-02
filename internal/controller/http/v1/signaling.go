package http_v1

import (
	"net/http"

	"rtc/internal/controller/websocket"

	gws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = gws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) handleSignaling(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	peerID := r.URL.Query().Get("id")
	h.logger.Info("New WebSocket connection", zap.String("peerID", peerID))
	peer := websocket.NewPeer(peerID, conn, h.hub)
	h.hub.Register(peer)

	go peer.ReadMessages()
	peer.WriteMessages()
}
