# ---- Build ----
FROM golang:alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o server .

# ---- Runtime ----
FROM alpine:3.22.2 AS runtime

RUN apk add --no-cache ca-certificates

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=build /app/server /app/

RUN chown app:app /app/server

USER app

EXPOSE 8080

CMD ["/app/server"]