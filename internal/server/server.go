package server

import (
	"net/http"

	http_v1 "rtc/internal/controller/http/v1"
	"rtc/internal/controller/websocket"

	"go.uber.org/zap"
)

func Start() error {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	hub := websocket.NewHub()
	go hub.Run()

	handler := http_v1.NewHandler(logger, hub)

	app := handler.InitRoutes()

	return http.ListenAndServe(":8080", app)
}
