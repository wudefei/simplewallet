FROM golang:1.23-alpine AS builder

ARG CGO=0
ENV GO111MODULE=on
ENV CGO_ENABLED=${CGO}
ENV GOOS=linux
ENV GOARCH=amd64


WORKDIR /build
COPY go.* ./
RUN go mod download

ARG env=dev
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build set -ex && \
    GOOS=${GOOS} GOARCH=${GOARCH} go build -o simplewallet main.go

FROM golang:1.23-alpine:buster

WORKDIR /work

COPY --from=builder /build/simplewallet .
COPY --from=builder /build/conf.yaml .

VOLUME ["/work/logs"]

EXPOSE 8080

CMD ["./simplewallet", "-conf=./conf.yaml"]



