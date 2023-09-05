#===============
# Stage 1: Build
#===============

FROM golang:1.20-alpine as builder

ENV BIN_REPO=github.com/distribyted/distribyted
ENV BIN_PATH=$GOPATH/src/$BIN_REPO

COPY . $BIN_PATH
WORKDIR $BIN_PATH

RUN apk add fuse-dev git gcc libc-dev g++ make

RUN BIN_OUTPUT=/bin/distribyted make build

#===============
# Stage 2: Run
#===============

FROM alpine:3

RUN apk add gcc libc-dev fuse-dev

COPY --from=builder /bin/distribyted /bin/distribyted
RUN chmod +x /bin/distribyted

RUN mkdir /distribyted-data

RUN echo "user_allow_other" >> /etc/fuse.conf
ENV DISTRIBYTED_FUSE_ALLOW_OTHER=true

ENTRYPOINT ["./bin/distribyted"]
