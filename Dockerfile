FROM golang:1.20 as builder

WORKDIR /api

COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o file-server.bin main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /api/file-server.bin /file-server

ENTRYPOINT ["/file-server"]