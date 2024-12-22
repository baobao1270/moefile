#!/usr/bin/env sh
set -e
exec /app/${APP_NAME} \
	-level "${LEVEL}" \
	-listen "${LISTEN}" \
	-origins "${ORIGINS}" \
	-proxies "${PROXIES}" \
	-root "${ROOT}" \
	-server "${SERVER}" \
	-xmltab "${XMLTAB}"
