// +build tools

package tools // import "github.com/AlekSi/alice/tools"

// direct dependencies
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/quasilyte/go-consistent"
	_ "github.com/reviewdog/reviewdog/cmd/reviewdog"
	_ "mvdan.cc/gofumpt/gofumports"
)

//go:generate go build -o ../bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go build -o ../bin/go-consistent github.com/quasilyte/go-consistent
//go:generate go build -o ../bin/reviewdog github.com/reviewdog/reviewdog/cmd/reviewdog
//go:generate go build -o ../bin/gofumports mvdan.cc/gofumpt/gofumports
