FROM golang:1.20.2-buster as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o capy-converter .

FROM alpine:latest
RUN apk add ffmpeg
COPY --from=builder /app/capy-converter .
CMD ["./capy-converter"]