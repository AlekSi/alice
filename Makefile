default: help

help:                                      ## Display this help message
	@echo "Please use \`make <target>\` where <target> is one of:"
	@grep '^[a-zA-Z]' $(MAKEFILE_LIST) | \
		awk -F ':.*?## ' 'NF==2 {printf "  %-26s%s\n", $$1, $$2}'

init:                                      ## Install development tools
	cd tools && go generate -x -tags=tools

format:                                    ## Format source code
	bin/gofumports -local github.com/AlekSi/alice -l -w .

test:                                      ## Run tests
	go install -v ./...
	go test -v -coverprofile=cover.out -covermode=count ./...
