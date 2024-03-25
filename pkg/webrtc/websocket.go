package webrtc

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: time.Second * 2,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
)

// ServeWs upgrades an HTTP request to Websocket, handles SFU signaling, and manages PeerConnections.
func (wr *SFU) ServeWs(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP request to Websocket
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	channelIDStr := r.URL.Query().Get("channel_id")
	if channelIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := &threadSafeWriter{
		Conn:  unsafeConn,
		Mutex: sync.Mutex{},
	}

	// When this frame returns close the Websocket
	defer c.Close() //nolint

	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{
		webrtc.RTPCodecTypeVideo,
		webrtc.RTPCodecTypeAudio,
	} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Add our new PeerConnection to global list
	wr.Lock()
	if _, ok := wr.peerConnectionsMap[channelID]; !ok {
		wr.peerConnectionsMap[channelID] = make(map[*peerConnectionState]bool)
	}

	peerConnectionState := &peerConnectionState{
		PeerConnection: peerConnection,
		websocket:      c,
	}
	wr.peerConnectionsMap[channelID][peerConnectionState] = true
	wr.Unlock()

	wr.setup(channelID, peerConnectionState)
	// Signal for the new PeerConnection
	wr.signalPeerConnections(channelID)

	message := &websocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			slog.Error("unable to read message: ", err)
			return
		}

		if err := json.Unmarshal(raw, &message); err != nil {
			slog.Error("unable to unmarshal message: ", err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				slog.Error("unable to unmarshal message: ", err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				slog.Error("unable to add ICE candidate: ", err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				slog.Error("unable to unmarshal message: ", err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				slog.Error("unable to set remote description: ", err)
				return
			}
		}
	}
}

func (wr *SFU) setup(channelID int64, peerConnection *peerConnectionState) {
	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		if writeErr := peerConnection.websocket.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			log.Println(writeErr)
		}
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			wr.signalPeerConnections(channelID)
		default:
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := wr.addTrack(channelID, t)
		defer wr.removeTrack(channelID, trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})
}
