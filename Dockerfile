FROM golang:1.23.1-alpine as build
RUN apk add build-base

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bot .

# Run stage
FROM alpine

WORKDIR /app
COPY --from=build /app/bot /app/ci6ndex

CMD ["/app/ci6ndex"]