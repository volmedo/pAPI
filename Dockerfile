# build
FROM golang:1.12.5 as builder
RUN adduser --disabled-password --gecos "" papiuser
WORKDIR /papi
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pkg/ pkg/
COPY cmd/ cmd/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /papi/srv ./cmd/server

# deploy
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /papi/srv /papi/srv
COPY --from=builder /papi/pkg/service/migrations /papi/migrations
USER papiuser
ENTRYPOINT ["/papi/srv", "-migrations=/papi/migrations"]
CMD ["-port=8080"]