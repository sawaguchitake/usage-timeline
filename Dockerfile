FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags='-s -w' -o webui ./cmd/webui

FROM gcr.io/distroless/static:nonroot
WORKDIR /root/
COPY --from=builder /app/webui .
COPY --from=builder /app/static ./static
EXPOSE 8080
USER nonroot
CMD ["./webui"]