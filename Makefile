MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash -o pipefail

export GO111MODULE = on

default: lint test integration

lint:
	# --enable=unparam
	golangci-lint run \
		--disable-all \
		--exclude-use-default=false \
		--exclude=package-comments \
		--enable=errcheck \
		--enable=goimports \
		--enable=ineffassign \
		--enable=revive \
		--enable=unused \
		--enable=structcheck \
		--enable=staticcheck \
		--enable=varcheck \
		--enable=deadcode \
		--enable=unconvert \
		--enable=misspell \
		--enable=prealloc \
		--enable=nakedret \
		--enable=typecheck \
		./...
test:
	go test ./... -race -cover -covermode=atomic -coverprofile=unit_coverage.cov

	@ echo "Converting the coverage file..."
	gocov convert ./unit_coverage.cov | gocov-xml > ./coverage.xml
integration:
	go run examples/tests/main.go test 
	go run examples/suites/main.go test 

sec:
	gosec -quiet ./...
