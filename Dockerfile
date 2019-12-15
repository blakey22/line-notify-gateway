FROM golang:1.13-alpine AS build

WORKDIR /app
RUN apk -u add git
COPY go.mod go.sum ./

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


FROM alpine:3.10
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY ./templates .
COPY --from=build /app .

EXPOSE 8080

CMD ["./server"]
