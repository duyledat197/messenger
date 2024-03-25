package webrtc

import "github.com/pion/webrtc/v3"

// Add to list of tracks and fire renegotation for all PeerConnections
func (wr *WebRTC) addTrack(channelID int64, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	wr.Lock()
	defer func() {
		wr.Unlock()
		wr.signalPeerConnections(channelID)
	}()

	// Create a new TrackLocal with the same codec as our incoming
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	wr.trackLocals[t.ID()] = trackLocal
	return trackLocal
}

// Remove from list of tracks and fire renegotation for all PeerConnections
func (wr *WebRTC) removeTrack(channelID int64, t *webrtc.TrackLocalStaticRTP) {
	wr.Lock()
	defer func() {
		wr.Unlock()
		wr.signalPeerConnections(channelID)
	}()

	delete(wr.trackLocals, t.ID())
}
