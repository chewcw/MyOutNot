FROM golang:1.17-alpine AS BUILD

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/exe
RUN apk add tzdata

FROM scratch
COPY --from=BUILD /app/bin/exe /usr/bin/exe
COPY --from=BUILD /etc/ssl/certs /etc/ssl/certs
COPY --from=BUILD /usr/share/zoneinfo /usr/share/zoneinfo
ENV GIN_MODE=release
ENV PORT=80
EXPOSE 80
ENTRYPOINT ["/usr/bin/exe"]