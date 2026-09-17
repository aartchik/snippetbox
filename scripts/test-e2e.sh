#!/bin/sh
set -eu

export APP_PORT=14000
export POSTGRES_PORT=15432
export REDIS_PORT=16380
export E2E_BASE_URL=http://127.0.0.1:14000

compose() {
    docker compose --project-name snippetbox-e2e "$@"
}

cleanup() {
    status=$?
    if [ "$status" -ne 0 ]; then
        compose logs
    fi
    compose down -v
    exit "$status"
}

trap cleanup EXIT INT TERM

compose up -d --build db redis app

attempt=1
while [ "$attempt" -le 60 ]; do
    if curl --fail --silent "$E2E_BASE_URL/ping" >/dev/null; then
        npm run test:e2e
        exit 0
    fi
    attempt=$((attempt + 1))
    sleep 2
done

echo "E2E application did not become ready" >&2
exit 1
