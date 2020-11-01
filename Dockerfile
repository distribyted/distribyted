#===============
# Stage 1: Build
#===============

FROM golang:1.15-alpine as builder

ENV BIN_REPO=github.com/distribyted/distribyted
ENV BIN_PATH=$GOPATH/src/$BIN_REPO

COPY . $BIN_PATH
WORKDIR $BIN_PATH

RUN apk add --update fuse-dev git gcc libc-dev g++

RUN go build -o /bin/distribyted -tags "release" cmd/distribyted/main.go

RUN ls /bin
#===============
# Stage 2: Run
#===============

FROM alpine:3

RUN apk add --update gcc libc-dev

COPY --from=builder /bin/distribyted /bin/distribyted

RUN chmod +x /bin/distribyted

ENTRYPOINT ["./bin/distribyted"]