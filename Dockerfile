FROM golang:latest AS build-env
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-handler

FROM alpine:latest
WORKDIR /app
COPY --from=build-env /app/notification-handler /app/notification-handler
CMD ["/app/notification-handler"]
