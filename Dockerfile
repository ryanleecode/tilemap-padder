
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go build ./tilemap-padder.go
RUN mv tilemap-padder /go/bin/app

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
