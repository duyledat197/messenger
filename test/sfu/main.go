package main

import (
	"context"
	"net/http"

	"openmyth/messgener/pkg/webrtc"
)

func main() {
	ctx := context.Background()
	sfu := webrtc.NewSFU()
	go sfu.Start(ctx)

	http.HandleFunc("/websocket", sfu.ServeWs)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
