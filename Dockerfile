# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial 

ENV p $GOPATH/src/github.com/labbsr0x/bindman-azure-dns-manager

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN GIT_COMMIT=$(git rev-parse --short HEAD 2> /dev/null || true) \
 && BUILDTIME=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ') \
 && VERSION=$(git describe --abbrev=0 --tags 2> /dev/null || true) \
 && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.Version=${VERSION:-unknow-version} \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.GitCommit=${GIT_COMMIT} \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.BuildTime=${BUILDTIME}" \
    -a -installsuffix cgo -o /bindman-azure-dns-manager src/main.go

# PKG
FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bindman-azure-dns-manager /go/bin/

VOLUME [ "/data" ]
ENTRYPOINT [ "/go/bin/bindman-azure-dns-manager" ]

CMD [ "serve" ]
