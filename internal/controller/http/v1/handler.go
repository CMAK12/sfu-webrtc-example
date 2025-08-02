package http_v1

import (
	"net/http"
	"rtc/internal/controller/websocket"

	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	hub    *websocket.Hub
}

func NewHandler(logger *zap.Logger, hub *websocket.Hub) *Handler {
	return &Handler{
		logger: logger,
		hub:    hub,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	hndl := h.loggingMiddleware(http.HandlerFunc(h.handleSignaling))

	mux.Handle("/ws", hndl)

	return mux
}
