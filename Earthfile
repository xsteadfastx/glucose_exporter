VERSION 0.7
ARG --global GO_VERSION=1.22.1
ARG --global GOLANGCILINT_VERSION=v1.54.2
ARG --global GORELEASER_VERSION=v1.24.0
FROM golang:$GO_VERSION-alpine3.18
WORKDIR /glucose_exporter

SAVE_CODE:
	COMMAND
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum
	SAVE ARTIFACT api AS LOCAL api
	SAVE ARTIFACT cmd AS LOCAL cmd
	SAVE ARTIFACT httpslog AS LOCAL httpslog
	SAVE ARTIFACT internal AS LOCAL internal

code:
	COPY go.mod go.sum ./
	COPY --dir api ./
	COPY --dir cmd ./
	COPY --dir httpslog ./
	COPY --dir internal ./
	COPY --dir vendor ./
	SAVE ARTIFACT . code

ci:
	FROM +deps
	ARG SNAPSHOT=true
	ARG GITHUB_TOKEN
	BUILD +test
	BUILD +lint
	COPY (+goreleaser/dist/metadata.json --SNAPSHOT=$SNAPSHOT --GITHUB_TOKEN=$GITHUB_TOKEN) /tmp/metadata.json
	ARG DOCKER_TAG=$(cat /tmp/metadata.json|jq -r .version)
	BUILD +docker --DOCKER_TAG=$DOCKER_TAG

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
	RUN apk add --no-cache \
		git \
		jq

lint:
	FROM +deps
	COPY +code/code ./
	COPY .golangci.yml .golangci.yml
	RUN golangci-lint run

goreleaser:
	FROM +deps
	ARG SNAPSHOT=true
	ARG GITHUB_TOKEN
	IF [ "$SNAPSHOT" = "true" ]
		ARG CMD="release --snapshot --clean"
	ELSE
		ARG CMD="release --clean"
	END
	COPY . .
	RUN goreleaser $CMD
	SAVE ARTIFACT dist dist

ssl-certs:
  RUN set -ex \
    && apk add --no-cache ca-certificates
  SAVE ARTIFACT /etc/ssl/certs/ca-certificates.crt ca-certificates.crt

docker:
	FROM scratch
	ARG DOCKER_TAG
  	COPY +ssl-certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
	COPY +goreleaser/dist/glucose_exporter_linux_amd64_v1/glucose_exporter /glucose_exporter
	VOLUME /var/cache/glucose_exporter
  	ENTRYPOINT ["/glucose_exporter"]
	CMD ["serve"]
	SAVE IMAGE --push ghcr.io/xsteadfastx/glucose_exporter:$DOCKER_TAG
