FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o restapi ./cmd/restapi/main.go

FROM scratch

COPY --from=builder /build/restapi /restapi

ENTRYPOINT ["/restapi"]