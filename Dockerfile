# syntax = docker/dockerfile:experimental
# Build Container
FROM golang:1.13.5 as data-sweeper-builder

ENV GO111MODULE on
WORKDIR /go/src/bitbucket.org/latonaio

COPY go.mod .

RUN go mod download

COPY . .

RUN go build


# Runtime Container
FROM alpine:3.12

RUN apk add --no-cache libc6-compat

ENV SERVICE=service-broker \
    POSITION=BackendService \
    AION_HOME="/var/lib/aion" \
    APP_DIR="${AION_HOME}/${POSITION}/${SERVICE}"

WORKDIR ${AION_HOME}

COPY --from=data-sweeper-builder /go/src/bitbucket.org/latonaio/data-sweeper-kube .

CMD ["./data-sweeper-kube"]
