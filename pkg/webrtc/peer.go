package webrtc

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"golang.org/x/sync/errgroup"
)

type peerConnectionState struct {
	UserID    string
	ChannelID int64
	*webrtc.PeerConnection
	websocket *threadSafeWriter
}

// Helper to make Gorilla Websockets threadsafe
type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

// signalPeerConnections updates each PeerConnection so that it is getting all the expected media tracks
func (wr *WebRTC) signalPeerConnections(channelID int64) {
	wr.Lock()
	defer func() {
		wr.Unlock()
		wr.dispatchKeyFrame()
	}()
	eg, _ := errgroup.WithContext(context.Background())

	attemptSync := func() (tryAgain bool) {
		for _, peerConnections := range wr.peerConnectionsMap {
			for peerConnection := range peerConnections {
				eg.Go(func() error {
					if peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
						wr.Lock()
						delete(wr.peerConnectionsMap[channelID], peerConnection)
						wr.Unlock()
						tryAgain = true // We modified the slice, start from the beginning
					}

					// map of sender we already are seanding, so we don't double send
					existingSenders := map[string]bool{}

					for _, sender := range peerConnection.GetSenders() {
						if sender.Track() == nil {
							continue
						}

						existingSenders[sender.Track().ID()] = true

						// If we have a RTPSender that doesn't map to a existing track remove and signal
						if _, ok := wr.trackLocals[sender.Track().ID()]; !ok {
							if err := peerConnection.RemoveTrack(sender); err != nil {
								tryAgain = true
							}
						}
					}

					// Don't receive videos we are sending, make sure we don't have loopback
					for _, receiver := range peerConnection.GetReceivers() {
						if receiver.Track() == nil {
							continue
						}

						existingSenders[receiver.Track().ID()] = true
					}

					// Add all track we aren't sending yet to the PeerConnection
					for trackID := range wr.trackLocals {
						if _, ok := existingSenders[trackID]; !ok {
							if _, err := peerConnection.AddTrack(wr.trackLocals[trackID]); err != nil {
								tryAgain = true
							}
						}
					}

					offer, err := peerConnection.CreateOffer(nil)
					if err != nil {
						tryAgain = true
					}

					if err := peerConnection.SetLocalDescription(offer); err != nil {
						tryAgain = true
					}

					offerString, err := json.Marshal(offer)
					if err != nil {
						tryAgain = true
					}

					if err = peerConnection.websocket.WriteJSON(&websocketMessage{
						Event: "offer",
						Data:  string(offerString),
					}); err != nil {
						tryAgain = true
					}

					return nil
				})
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				wr.signalPeerConnections(channelID)
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
