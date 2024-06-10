package httpclient

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type RoundTripper struct{}

func NewRoundTripper() http.RoundTripper {
	return &RoundTripper{}
}

func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	reqBytes, _ := jsoniter.Marshal(req)
	reqID := uuid.NewString()
	slog.Info("request: ", slog.String("request_id", reqID), slog.String("request", string(reqBytes)))
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	slog.Info("response: ", slog.String("request_id", reqID), slog.Int("status", resp.StatusCode), slog.Duration("duration", time.Since(startTime)))

	return resp, nil
}
