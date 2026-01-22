package service

import (
	"net/http"

	"{{.Computed.common_module_final}}/log"
	"github.com/gorilla/websocket"
)

const (
	wsReadBufferSize  = 1024
	wsWriteBufferSize = 1024
)

var upgrader = websocket.Upgrader{
	// Allow connections from any origin
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  wsReadBufferSize,
	WriteBufferSize: wsWriteBufferSize,
}

// Ws handles WebSocket connections at the /ws endpoint.
func (s *{{.Computed.service_name_capitalized}}Service) Ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("failed to upgrade connection to websocket")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	log.Info("websocket connection established")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Warn("websocket unexpected close")
			}
			break
		}

		log.WithField("message", string(message)).Debug("received websocket message")

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.WithError(err).Error("failed to write websocket message")
			break
		}
	}

	log.Info("websocket connection closed")
}

// WsWithHeartbeat demonstrates a more advanced WebSocket handler with heartbeat.
// Uncomment and add to NewWSHandler if needed.
/*
func (s *{{.Computed.service_name_capitalized}}Service) WsWithHeartbeat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("failed to upgrade connection to websocket")
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Set pong handler
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start ping ticker
	ticker := time.NewTicker(54 * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})

	// Read messages
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.WithError(err).Warn("websocket unexpected close")
				}
				return
			}
			log.WithField("message", string(message)).Debug("received websocket message")

			// Process message here
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.WithError(err).Error("failed to write websocket message")
				return
			}
		}
	}()

	// Send pings
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
*/
