FROM      --platform=$BUILDPLATFORM alpine:3.21 AS builder
RUN       apk add --no-cache tzdata ca-certificates && mkdir /data

FROM      alpine:3.21
ARG       TARGETPLATFORM
ARG       APP_NAME=moefile
COPY      --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY      --from=builder /usr/share/ca-certificates /usr/share/ca-certificates
COPY      --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY      --from=builder /data /data
COPY      bin/build/${TARGETPLATFORM}/${APP_NAME} /app/${APP_NAME}
COPY      scripts/container-init.sh /app/container-init.sh
WORKDIR   /app
EXPOSE    3328
VOLUME    /data
ENV       APP_NAME=${APP_NAME} \
          LEVEL=inf \
          LISTEN=:3328 \
          ORIGINS=* \
          PROXIES=127.0.0.1 \
          ROOT=/data \
          SERVER=${APP_NAME} \
          XMLTAB=true
CMD       ["/app/container-init.sh"]
