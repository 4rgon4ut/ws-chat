package transport

import (
	"bytes"
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
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
	conn   *websocket.Conn
	ReadCh chan []byte
	SendCh chan []byte
	exitCh chan struct{}
}

// NewWSAdapter ...
func NewWSAdapter(conn *websocket.Conn, rCh chan []byte, wCh chan []byte) *WSAdapter {
	return &WSAdapter{
		conn:   conn,
		ReadCh: rCh,
		SendCh: wCh,
	}
}

//
func (wa *WSAdapter) reading(exitCh chan *WSAdapter) {
	defer func() {
		exitCh <- wa
		wa.conn.Close()
	}()
	wa.conn.SetReadLimit(maxMessageSize)
	wa.conn.SetReadDeadline(time.Now().Add(pongWait))
	wa.conn.SetPongHandler(func(string) error { wa.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := wa.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		wa.ReadCh <- message
	}
}

func (wa *WSAdapter) writing() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wa.conn.Close()
	}()
	for {
		select {
		case message, ok := <-wa.SendCh:
			wa.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				wa.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := wa.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(wa.SendCh)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-wa.SendCh)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			wa.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wa.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//
// func (wa *WSAdapter) close() {

// }

// Serve ...
func (wa *WSAdapter) Serve(exitCh chan *WSAdapter) {
	go wa.reading(exitCh)
	go wa.writing()
}
