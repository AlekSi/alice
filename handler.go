package alice

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

// Printf is a log.Printf-like function that can be used for logging.
type Printf func(format string, a ...interface{})

// Responder is a function that should be implemented to handle Yandex.Dialogs requests.
//
// Passed context is derived from HTTP request's context with added handler's timeout.
// It is canceled when the request is canceled (see https://golang.org/pkg/net/http/#Request.Context)
// or on timeout.
//
// Only response payload can be returned; other response fields (session, version) will be set automatically.
// If error is returned, it is logged with error logger, and 500 Internal server error is sent in response.
type Responder func(ctx context.Context, request *Request) (*ResponsePayload, error)

// Handler accepts Yandex.Dialogs requests, decodes them, handles "ping" requests itself,
// and delegates other requests to responder.
type Handler struct {
	r Responder

	Timeout time.Duration // responder's timeout
	Errorf  Printf        // error logger

	// debugging options
	Debugf        Printf // debug logger
	Indent        bool   // indent requests and responses
	StrictDecoder bool   // disallow unexpected fields in requests
}

// NewHandler creates new handler with given responder and default timeout (3s).
// Exported fields of the returned object can be changed before usage.
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

func internalError(rw http.ResponseWriter) {
	http.Error(rw, "Internal server error.", 500)
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), h.Timeout)
	defer cancel()

	if h.Debugf != nil {
		if h.Indent {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				h.errorf("Failed to read request: %s.", err)
				internalError(rw)
				return
			}

			var body bytes.Buffer
			if err = json.Indent(&body, b, "", "  "); err != nil {
				h.errorf("Failed to indent request: %s.", err)
				internalError(rw)
				return
			}
			req.Body = ioutil.NopCloser(&body)
			req.ContentLength = int64(body.Len())
			req.TransferEncoding = nil
		}

		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			h.errorf("Failed to dump request: %s.", err)
			internalError(rw)
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
		internalError(rw)
		return
	}
	if payload == nil {
		h.errorf("Responder returned nil payload without error.")
		internalError(rw)
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

	var body bytes.Buffer
	encoder := json.NewEncoder(&body)
	if h.Indent {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(response); err != nil {
		h.errorf("Failed to encode response body: %s.", err)
		internalError(rw)
		return
	}
	h.debugf("Response body:\n%s", body.Bytes())

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Header().Set("Content-Length", strconv.Itoa(body.Len()))
	rw.WriteHeader(200)
	if _, err := rw.Write(body.Bytes()); err != nil {
		h.errorf("Failed to write response body: %s.", err)
	}
}

// check interface
var _ http.Handler = (*Handler)(nil)
