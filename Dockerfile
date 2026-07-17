FROM golang:1.25.2-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /server /server

RUN mkdir -p /uploads && chown app:app /uploads

USER app

EXPOSE 8000

CMD ["/server"]
