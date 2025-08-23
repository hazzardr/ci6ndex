FROM golang:1.25.0-alpine AS build
RUN apk add build-base

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o ./bin/civ .

# Run stage
FROM alpine

WORKDIR /app
COPY --from=build /app/bin/civ /app/ci6ndex

CMD ["/app/ci6ndex", "bot", "serve"]
