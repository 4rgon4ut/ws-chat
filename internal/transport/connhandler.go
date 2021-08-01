package transport

import (
	"fmt"
	"time"

	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// WSAdapter  ...
type WSAdapter struct {
	Conn   *websocket.Conn
	SendCh chan string
	exitCh chan *WSAdapter
}

// NewWSAdapter ...
func NewWSAdapter(conn *websocket.Conn, sendCh chan string, exitCh chan *WSAdapter) *WSAdapter {
	return &WSAdapter{
		Conn:   conn,
		SendCh: sendCh,
		exitCh: exitCh,
	}
}

// Listen ...
func (wa *WSAdapter) Listen() {
	defer func() {
		wa.exitCh <- wa
		wa.Conn.Close()
	}()

	for {
		messageType, message, err := wa.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("read error:", err)
			}
			log.Infof("connection with %s closed", wa.Conn.RemoteAddr())
			return // Calls the deferred function, i.e. closes the connection on error
		}
		if messageType == websocket.TextMessage {
			// Broadcast the received message
			wa.SendCh <- fmt.Sprintf("[%s]: %s", wa.Conn.RemoteAddr(), message)
		} else {
			log.Error("websocket message received of type", messageType)
		}
	}
}

func (wa *WSAdapter) Write(message string) {
	if err := wa.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Error("write error:", err)
		wa.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		wa.Conn.Close()
		wa.exitCh <- wa
	}
}
