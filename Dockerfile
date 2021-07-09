FROM golang:1.16.5-alpine3.13 as builder

ARG CI_DATE
ARG CI_COMMIT_REF_NAME
ARG CI_COMMIT_SHA

RUN apk add --no-cache ca-certificates

RUN env

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
        go build -ldflags="-X 'main.BuildTime=${CI_DATE}' -X 'main.BuildVersion=${CI_COMMIT_REF_NAME}' -X 'main.CommitId=${CI_COMMIT_SHA}'" -o docker-hub-limit-expoter

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt \
     /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/docker-hub-limit-expoter /docker-hub-limit-expoter
LABEL org.opencontainers.image.created=${CI_DATE}
LABEL org.opencontainers.image.authors="Konstantin Fomin <konst@mhn.lv>"
LABEL org.opencontainers.image.version=${CI_COMMIT_REF_NAME}
LABEL org.opencontainers.image.revision=${CI_COMMIT_SHA}
LABEL org.opencontainers.image.title="Docker hub limit exporter"
LABEL org.opencontainers.image.source="https://github.com/fomk/docker-hub-limit-exporter"

ENTRYPOINT ["/docker-hub-limit-expoter"]
