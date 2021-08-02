package wstools

import (
	"fmt"

	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

// Adaptor binds to single websocket connection and server it
// send and read and fromat messages
type Adaptor struct {
	Conn   *websocket.Conn
	SendCh chan string
	exitCh chan *Adaptor
}

// NewAdaptor ...
func NewAdaptor(conn *websocket.Conn, sendCh chan string, exitCh chan *Adaptor) *Adaptor {
	return &Adaptor{
		Conn:   conn,
		SendCh: sendCh,
		exitCh: exitCh,
	}
}

// Listen connection in a loop, brakes on error and send signal to hub unregister channel
// than close connection
func (ad *Adaptor) Listen() {
	defer func() {
		ad.exitCh <- ad
		ad.Conn.Close()
	}()

	for {
		messageType, message, err := ad.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("read error:", err)
			}
			log.Infof("connection with %s closed", ad.Conn.RemoteAddr())
			return // Calls the deferred function, i.e. closes the connection on error
		}
		if messageType == websocket.TextMessage {
			// Broadcast the received message
			ad.SendCh <- fmt.Sprintf("[%s]: %s", ad.Conn.RemoteAddr(), message)
		} else {
			log.Error("websocket message received of type", messageType)
		}
	}
}

// Write ...
func (ad *Adaptor) Write(message string) {
	if err := ad.Conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Error("write error:", err)
		ad.Conn.WriteMessage(websocket.CloseMessage, []byte{})
		ad.Conn.Close()
		ad.exitCh <- ad
	}
}
