FROM golang:1.22.3-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates bash curl

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["/root/main"]