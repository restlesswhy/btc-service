FROM golang:1.18 AS builder
WORKDIR /go/src/github.com/restlesswhy/btc-service/

#COPY go.mod ./
#COPY go.sum ./
#RUN go mod download

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/app/main.go

FROM alpine:latest AS app
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/restlesswhy/btc-service/app ./
ENTRYPOINT [ "/app/app" ]  
