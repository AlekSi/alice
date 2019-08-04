#!/bin/sh

set -eux

go test -v -short -coverprofile cover.out ./...

curl -s https://codecov.io/bash | bash -s -- -X fix -e GOLANG_VERSION
