package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"

	"openmyth/messgener/util"
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

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// TODO: for security just accept some origin
	CheckOrigin:       func(r *http.Request) bool { return true },
	HandshakeTimeout:  time.Second * 2,
	EnableCompression: true,
}

// Client is a middleman between the websocket connection and the engine.
type Client[T any] struct {
	ID        string
	ChannelID int64
	UserID    string

	Engine *Engine[T]

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan T

	Impl Impl[T]
}

// readPump pumps messages from the websocket connection to the engine.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client[T]) readPump() {
	defer func() {
		c.Engine.unregister <- c
		if err := c.Impl.UnRegister(c.UserID); err != nil {
			log.Fatalf("unable to unregister: %v", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Fatalf("unable to close client: %v", err)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		slog.Error("unable to read deadline: %v", err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			slog.Error("unable to set read deadline: %v", err)

			return err
		}

		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		var msg T
		if err := json.Unmarshal(message, msg); err != nil {
			log.Printf("unable to unmarshal message: %v", err)
			continue
		}

		if err := c.Impl.Execute(c, msg); err != nil {
			log.Printf("unable to impl read message: %v", err)
		}
	}
}

// writePump pumps messages from the engine to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client[T]) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("unable to set write deadline: %v", err)
			}
			if !ok {
				// The engine closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("unable to set write message: %v", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("unable to set next message: %v", err)
				return
			}

			// TODO: support multi marshaler
			bData, err := json.Marshal(message)
			if err != nil {
				log.Printf("unable to marshal: %v", err)
				return
			}
			if _, err := w.Write(bData); err != nil {
				log.Printf("unable to write message: %v", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Printf("unable to write message: %v", err)
				}
				// TODO: support multi marshaler
				bData, err := json.Marshal(<-c.Send)
				if err != nil {
					log.Printf("unable to marshal: %v", err)
					return
				}

				if _, err := w.Write(bData); err != nil {
					log.Printf("unable to write message: %v", err)
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Printf("unable to close writer: %v", err)

				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("unable to set write deadline: %v", err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("unable to set write message: %v", err)

				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs[T any](engine *Engine[T], w http.ResponseWriter, r *http.Request) {
	schema, token, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	if !ok || strings.ToLower(schema) != "bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	channelIDStr := r.URL.Query().Get("channel_id")
	if channelIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	channelID := cast.ToInt64(channelIDStr)

	userClaims, err := util.VerifyToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// TODO: set up response header
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	client := &Client[T]{
		ID:        uuid.NewString(),
		UserID:    userClaims.Id,
		ChannelID: int64(channelID),
		Engine:    engine,
		conn:      conn,
		Send:      make(chan T),
	}

	// Allow collection of memory referenced by the caller by doing all work in
	client.Engine.register <- client

	if err := client.Impl.Register(client.ID, client.UserID); err != nil {
		log.Printf("unable to register: %v", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
