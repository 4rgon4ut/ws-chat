package wstools

import (
	"encoding/json"
	"fmt"

	"github.com/bestpilotingalaxy/ws-chat/internal/packets/jrpc"
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
)

// Adaptor binds to single websocket connection and server it
// send and read and fromat messages
type Adaptor struct {
	Conn   *websocket.Conn
	SendCh chan string
	jrpcCh chan *jrpc.Request
	exitCh chan *Adaptor
}

// NewAdaptor ...
func NewAdaptor(conn *websocket.Conn, sendCh chan string, jrpcCh chan *jrpc.Request, exitCh chan *Adaptor) *Adaptor {
	return &Adaptor{
		Conn:   conn,
		SendCh: sendCh,
		jrpcCh: jrpcCh,
		exitCh: exitCh,
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

		} else if messageType == websocket.BinaryMessage {
			jrpcCall := &jrpc.Request{}

			if err := json.Unmarshal(message, jrpcCall); err != nil {
				go ad.Write(fmt.Sprint("Wrong JRPC call format: ", err))
				continue
			}
			ad.jrpcCh <- jrpcCall

		} else {
			log.Error("websocket message received of type", messageType)
		}
	}
}
