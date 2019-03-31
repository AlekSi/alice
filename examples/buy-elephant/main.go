// +build ignore

// Port of https://github.com/yandex/alice-skills/tree/master/python/buy-elephant
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/AlekSi/alice"
)

//nolint:gochecknoglobals
var (
	sessionM       sync.Mutex
	sessionStorage = make(map[string][]string)
)

func getSuggests(userID string) []alice.ResponseButton {
	res := make([]alice.ResponseButton, 0, 2)

	// select two first suggestions
	for _, suggest := range sessionStorage[userID] {
		res = append(res, alice.ResponseButton{
			Title: suggest,
			Hide:  true,
		})
	}
	if len(res) > 2 {
		res = res[:2]
	}

	// remove first stored suggestion
	if len(sessionStorage[userID]) != 0 {
		sessionStorage[userID] = sessionStorage[userID][1:]
	}

	// add Yandex.Market suggestion
	if len(res) < 2 {
		res = append(res, alice.ResponseButton{
			Title: "Ладно",
			URL:   "https://market.yandex.ru/search?text=слон",
			Hide:  true,
		})
	}

	return res
}

func main() {
	flag.Parse()

	h := alice.NewHandler(func(ctx context.Context, request *alice.Request) (*alice.ResponsePayload, error) {
		sessionM.Lock()
		defer sessionM.Unlock()

		userID := request.Session.UserID
		if request.Session.New {
			sessionStorage[userID] = []string{
				"Не хочу.",
				"Не буду.",
				"Отстань!",
			}

			return &alice.ResponsePayload{
				Text:    "Привет! Купи слона!",
				Buttons: getSuggests(userID),
			}, nil
		}

		req := strings.ToLower(request.Request.OriginalUtterance)
		for _, expected := range []string{"ладно", "куплю", "покупаю", "хорошо"} {
			if req == expected {
				return &alice.ResponsePayload{
					Text: "Слона можно найти на Яндекс.Маркете!",
				}, nil
			}
		}

		return &alice.ResponsePayload{
			Text:    fmt.Sprintf("Все говорят %q, а ты купи слона!", req),
			Buttons: getSuggests(userID),
		}, nil
	})

	h.Errorf = log.New(os.Stdout, "error: ", 0).Printf
	h.Debugf = log.New(os.Stdout, "debug: ", 0).Printf
	h.Indent = true
	http.Handle("/", h)

	const addr = "127.0.0.1:8080"
	log.Printf("Listening on http://%s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
