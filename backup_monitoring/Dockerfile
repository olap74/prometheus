FROM golang:1.14 as builder

ENV USER=appuser
ENV UID=10001

RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /app
COPY app .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -ldflags="-w -s" -o /app/prom_exporter

FROM scratch

COPY --from=builder /app/prom_exporter /prom_exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER appuser:appuser

ENTRYPOINT ["/prom_exporter"]
