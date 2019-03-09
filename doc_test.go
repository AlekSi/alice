package alice_test

import (
	"context"
	"log"
	"net/http"

	"github.com/AlekSi/alice"
)

func Example() {
	h := alice.NewHandler(func(ctx context.Context, request *alice.Request) (*alice.ResponsePayload, error) {
		return &alice.ResponsePayload{
			Text:       "Bye!",
			EndSession: true,
		}, nil
	})

	h.Errorf = log.Printf
	http.Handle("/", h)

	const addr = "127.0.0.1:8080"
	log.Printf("Listening on http://%s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
