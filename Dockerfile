# BUILD
FROM golang:1.13-alpine3.10 as builder

RUN apk --no-cache --no-progress add git

WORKDIR $GOPATH/src/github.com/labbsr0x/bindman-azure-dns-manager

COPY go.mod go.sum ./
ENV GO111MODULE=on
RUN go mod download

COPY . .

RUN GIT_COMMIT=$(git rev-parse --short HEAD 2> /dev/null || true) \
 && BUILDTIME=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ') \
 && VERSION=$(git describe --abbrev=0 --tags 2> /dev/null || true) \
 && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.Version=${VERSION:-unknow-version} \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.GitCommit=${GIT_COMMIT} \
    -X github.com/labbsr0x/bindman-azure-dns-manager/src/version.BuildTime=${BUILDTIME}" \
    -a -installsuffix cgo -o /bindman-azure-dns-manager src/main.go

# PKG
FROM alpine:3.10

RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

COPY --from=builder /bindman-azure-dns-manager /go/bin/

VOLUME [ "/data" ]
ENTRYPOINT [ "/go/bin/bindman-azure-dns-manager" ]

CMD [ "serve" ]
