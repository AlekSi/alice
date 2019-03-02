package alice // import "github.com/AlekSi/alice"

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

type Responder func(ctx context.Context, request *Request, response *Response) error

type Handler struct {
	r      Responder
	Errorf Printf
	Debugf Printf
}

func NewHandler(r Responder) *Handler {
	return &Handler{
		r: r,
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

func pingResponder(ctx context.Context, request *Request, response *Response) error {
	response.Response.Text = "pong"
	response.Response.EndSession = true
	return nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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
	if h.Debugf != nil {
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

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()
	response := NewResponse(request)
	if err := r(ctx, request, response); err != nil {
		h.errorf("Responder failed: %s.", err)
		http.Error(rw, "Internal server error.", 500)
		return
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if h.Debugf != nil {
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
