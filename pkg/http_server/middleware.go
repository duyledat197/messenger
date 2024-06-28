package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	mtdt "openmyth/messgener/pkg/metadata"
	"openmyth/messgener/util"
)

type middlewareFunc func(http.Handler) http.Handler

var (
	// Authorization is the name of the authorization header
	Authorization = "Authorization"
	// Bearer is the name of the bearer token
	Bearer = "Bearer"

	ignoredAPIs         []string
	invalidateCacheAPIs []string
	ignoredForLogAPIs   []string
)

type payloadKeys struct{}

// GetClientIP get client IP from HTTP request
func GetClientIP(req *http.Request) string {
	md, ok := metadata.FromIncomingContext(req.Context())
	if !ok {
		return ""
	}
	clientIP := md.Get(mtdt.MDXForwardedFor)
	if len(clientIP) == 0 {
		return ""
	}

	return clientIP[0]
}

// allowCORS sets up CORS headers for the HTTP handler.
//
// It takes an http.Handler as a parameter and returns an http.Handler.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	})
}

// MapMetaDataFunc defines a function that extracts metadata from the request
type MapMetaDataFunc func(context.Context, *http.Request) metadata.MD

// MapMetaDataWithBearerToken defines a function that extracts authorization information from the request header
// and creates metadata based on the bearer token found in the authorization header.
func MapMetaDataWithBearerToken() MapMetaDataFunc {
	return func(ctx context.Context, r *http.Request) metadata.MD {
		md := metadata.MD{}
		authorization := r.Header.Get(Authorization)
		if authorization != "" {
			schema, bearerToken, ok := strings.Cut(authorization, " ")
			if !ok || strings.ToLower(schema) != strings.ToLower(Bearer) {
				return md
			}
			payload, err := util.VerifyToken(bearerToken)
			if err == nil {
				md = metadata.Join(md, mtdt.ImportUserInfoToCtx(&mtdt.Payload{
					UserID: payload.Id,
					Ip:     GetClientIP(r),
					Token:  bearerToken,
				}))
			}
		}

		return md
	}
}

// func MapMetaDataWithBearerToken(authenticator authenticate.Authenticator) mapMetaDataFunc {
// 	return func(ctx context.Context, r *http.Request) metadata.MD {
// 		md := mtdt.ImportIpToCtx(GetClientIP(r))

// 		authorization := r.Header.Get(Authorization)

// 		if authorization != "" {
// 			bearerToken := strings.Split(authorization, Bearer+" ")
// 			if len(bearerToken) < 2 {
// 				return md
// 			}
// 			token := bearerToken[1]
// 			payload, err := authenticator.Verify(token)
// 			if err != nil {
// 				return md
// 			}
// 			payload.Token = token

// 			md = metadata.Join(md, mtdt.ImportUserInfoToCtx(payload))
// 		}

// 		return md
// 	}
// }

// Response is the response struct for the http server
type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Data    any      `json:"data"`
}

// ErrorResponse generates an error response with the provided code and error message.
func ErrorResponse(w http.ResponseWriter, code int, err error) {
	resp := &Response{
		Code:    code,
		Message: err.Error(),
		Details: []string{},
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// DataResponse generates a data response with the provided data and writes it to the http.ResponseWriter.
func DataResponse(w http.ResponseWriter, data any) {
	resp := &Response{
		Data: data,
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// forwardErrorResponse forwards an error response to the client based on the provided error.
func forwardErrorResponse(ctx context.Context, s *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	sta := status.Convert(err)
	errStr := sta.Message()
	firstColonPos := strings.Index(errStr, ":")

	if firstColonPos > 0 {
		errStr = errStr[:firstColonPos]
	}

	runtime.DefaultHTTPErrorHandler(ctx, s, m, w, r, errors.New(errStr))
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// NewWrapResponseWriter creates a new wrapped response writer based on the provided http.ResponseWriter.
//
// It takes an http.ResponseWriter as a parameter and returns a pointer to a responseWriterWrapper.
func NewWrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK}
}

// WriteHeader sets the status code for the response writer.
//
// It takes an integer code as a parameter.
func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
