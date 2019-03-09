package alice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type Printf func(format string, a ...interface{})

type Responder func(ctx context.Context, request *Request) (*ResponsePayload, error)

type Handler struct {
	r              Responder
	Timeout        time.Duration
	Errorf         Printf
	Debugf         Printf
	IndentResponse bool
	StrictDecoder  bool
}

func NewHandler(r Responder) *Handler {
	return &Handler{
		r:       r,
		Timeout: 3 * time.Second,
	}
}

func (h *Handler) errorf(format string, a ...interface{}) {
	if h.Errorf != nil {
		h.Errorf(format, a...)
	}
}

func (h *Handler) debugf(format string, a ...interface{}) {
	if h.Debugf != nil {
		h.Debugf(format, a...)
	}
}

func pingResponder(ctx context.Context, request *Request) (*ResponsePayload, error) {
	return &ResponsePayload{
		Text:       "pong",
		EndSession: true,
	}, nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), h.Timeout)
	defer cancel()

	if h.Debugf != nil {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			h.errorf("Failed to dump request: %s.", err)
			http.Error(rw, "Internal server error.", 500)
			return
		}
		h.debugf("Request:\n%s", b)
	}

	request := new(Request)
	decoder := json.NewDecoder(req.Body)
	if h.StrictDecoder {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(request); err != nil {
		h.errorf("Failed to read or decode request body: %s.", err)
		http.Error(rw, "Failed to decode request body.", 400)
		return
	}

	r := h.r
	if request.Request.Type == SimpleUtterance && request.Request.OriginalUtterance == "ping" {
		r = pingResponder
	}

	payload, err := r(ctx, request)
	if err != nil {
		h.errorf("Responder failed: %s.", err)
		http.Error(rw, "Internal server error.", 500)
		return
	}
	if payload == nil {
		h.errorf("Responder returned nil payload without error.")
		http.Error(rw, "Internal server error.", 500)
		return
	}
	response := &Response{
		Response: *payload,
		Session: ResponseSession{
			SessionID: request.Session.SessionID,
			MessageID: request.Session.MessageID,
			UserID:    request.Session.UserID,
		},
		Version: request.Version,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if h.IndentResponse {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(response); err != nil {
		h.errorf("Failed to encode response body: %s.", err)
		http.Error(rw, "Internal server error.", 500)
		return
	}
	h.debugf("Response:\n%s", buf.Bytes())

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	rw.WriteHeader(200)
	if _, err := rw.Write(buf.Bytes()); err != nil {
		h.errorf("Failed to write response body: %s.", err)
	}
}

// check interface
var _ http.Handler = (*Handler)(nil)
