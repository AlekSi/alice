# alice

[![Actions](https://img.shields.io/badge/tested%20with-actions-success.svg?logo=github)](https://github.com/AlekSi/alice/actions)
[![Codecov](https://codecov.io/gh/AlekSi/alice/branch/master/graph/badge.svg)](https://codecov.io/gh/AlekSi/alice)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlekSi/alice)](https://goreportcard.com/report/github.com/AlekSi/alice)
[![GoDoc](https://godoc.org/github.com/AlekSi/alice?status.svg)](https://godoc.org/github.com/AlekSi/alice)

Package alice provides helpers for developing skills for Alice virtual assistant
via Yandex.Dialogs platform.

# Example

```go
responder := func(ctx context.Context, request *alice.Request) (*alice.ResponsePayload, error) {
    return &alice.ResponsePayload{
        Text:       "Bye!",
        EndSession: true,
    }, nil
}

h := alice.NewHandler(responder)
h.Errorf = log.Printf
http.Handle("/", h)

const addr = "127.0.0.1:8080"
log.Printf("Listening on http://%s ...", addr)
log.Fatal(http.ListenAndServe(addr, nil))
```

See [documentation](https://godoc.org/github.com/AlekSi/alice) and [examples](examples).

# License

Copyright (c) 2019 Alexey Palazhchenko. [MIT-style license](LICENSE).

Tests use the following resources:
* https://freesound.org/people/prucanada/sounds/415341/ by prucanada.
