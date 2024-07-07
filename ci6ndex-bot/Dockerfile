# Build stage
FROM golang:1.22 as build

WORKDIR /app
ADD .. /app

RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ci6ndex ./main.go

# Run stage
FROM alpine

WORKDIR /app
COPY --from=build /app/ci6ndex /app/ci6ndex
COPY --from=build /app/templates /app/templates

CMD ["/app/ci6ndex", "discord", "start"]