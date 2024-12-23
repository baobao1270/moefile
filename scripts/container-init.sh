#!/usr/bin/env sh
set -e
XMLTAB_ARGS=""
if [ "${XMLTAB}" == "true" ]; then
	XMLTAB_ARGS="-xmltab"
fi
exec /app/${APP_NAME} \
	-level "${LEVEL}" \
	-listen "${LISTEN}" \
	-origins "${ORIGINS}" \
	-proxies "${PROXIES}" \
	-root "${ROOT}" \
	-server "${SERVER}" \
	${XMLTAB_ARGS}
