VERSION 0.7
ARG --global GO_VERSION=1.22.0
ARG --global GOLANGCILINT_VERSION=v1.54.2
ARG --global GORELEASER_VERSION=v1.24.0
FROM golang:$GO_VERSION-alpine3.18
WORKDIR /glucose_exporter

SAVE_CODE:
	COMMAND
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum
	SAVE ARTIFACT api AS LOCAL api
	SAVE ARTIFACT httpslog AS LOCAL httpslog
	SAVE ARTIFACT internal AS LOCAL internal

code:
	COPY go.mod go.sum ./
	COPY --dir api ./
	COPY --dir httpslog ./
	COPY --dir internal ./
	COPY --dir vendor ./
	SAVE ARTIFACT . code

ci:
	BUILD +test
	BUILD +lint

test:
	COPY +code/code ./
	RUN go test -v ./...

tidy:
	RUN apk add --no-cache git
	COPY +code/code ./
	RUN go mod tidy
	RUN go mod vendor
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum
	SAVE ARTIFACT vendor AS LOCAL vendor

deps:
	FROM +base
	RUN (cd /tmp; go install github.com/golangci/golangci-lint/cmd/golangci-lint@$GOLANGCILINT_VERSION)
	RUN (cd /tmp; go install github.com/goreleaser/goreleaser@$GORELEASER_VERSION)

lint:
	FROM +deps
	COPY +code/code ./
	COPY .golangci.yml .golangci.yml
	RUN golangci-lint run
