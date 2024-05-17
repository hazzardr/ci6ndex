# Build stage
FROM golang:1.22 as build

WORKDIR /app
ADD . /app

RUN go mod download
RUN go build -o ci6ndex ./main.go

# Final stage
FROM golang:1.22-alpine

WORKDIR /app
COPY --from=build /app/ci6ndex /app/ci6ndex

CMD ["/app/ci6ndex", "discord", "start"]