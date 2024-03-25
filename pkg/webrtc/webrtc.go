package webrtc

import (
	"context"
	"sync"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"golang.org/x/sync/errgroup"
)

type WebRTC struct {
	peerConnectionsMap map[int64]map[*peerConnectionState]bool
	trackLocals        map[string]*webrtc.TrackLocalStaticRTP

	isActive bool
	sync.Mutex
}

// NewWebRTC initializes a new WebRTC instance.
func NewWebRTC() *WebRTC {
	return &WebRTC{
		peerConnectionsMap: make(map[int64]map[*peerConnectionState]bool),
		trackLocals:        make(map[string]*webrtc.TrackLocalStaticRTP),
		isActive:           true,
		Mutex:              sync.Mutex{},
	}
}
func (wr *WebRTC) Start(_ context.Context) error {

	// request a keyframe every 3 seconds
	for range time.NewTicker(time.Second * 3).C {
		if wr.isActive {
			wr.dispatchKeyFrame()
		}
	}

	return nil
}

// Stop stops the WebRTC.
func (wr *WebRTC) Stop(_ context.Context) error {
	wr.isActive = false

	return nil
}

func (wr *WebRTC) dispatchKeyFrame() {
	wr.Lock()
	defer wr.Unlock()

	eg, _ := errgroup.WithContext(context.Background())

	for _, peerConnections := range wr.peerConnectionsMap {
		for peerConnection := range peerConnections {
			for _, receiver := range peerConnection.GetReceivers() {
				eg.Go(func() error {
					if receiver.Track() == nil {
						return nil
					}

					if err := peerConnection.WriteRTCP([]rtcp.Packet{
						&rtcp.PictureLossIndication{
							MediaSSRC: uint32(receiver.Track().SSRC()),
						},
					}); err != nil {
						return nil
					}

					return nil
				})
			}
		}
	}

	eg.Wait()
}
