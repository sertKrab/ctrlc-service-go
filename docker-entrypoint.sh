#!/bin/sh
set -e

# The API needs the schema migrated before it starts. migrate connects to the
# database itself, so we do not need pg_isready/nc (absent from the alpine
# runtime image) — we just retry migrate until the DB accepts connections.
#
# When this container is orchestrated with a readiness gate (docker compose
# `depends_on: condition: service_healthy`, k8s readiness probes) the DB is
# already up on the first try. Under plain `docker run`, a compose without
# health gating, or a managed DB that briefly refuses connections, the first
# attempt can lose the race — the retry keeps the container from crash-looping.
# A genuine migration error still surfaces after the attempt budget is spent.
MIGRATE_RETRIES="${MIGRATE_RETRIES:-30}"
MIGRATE_RETRY_DELAY="${MIGRATE_RETRY_DELAY:-2}"

echo "[entrypoint] running migrations (up to ${MIGRATE_RETRIES} attempts)..."
i=1
while true; do
  # `if cmd` disables set -e for cmd, so a failing migrate does not exit here.
  if ./migrate up; then
    break
  fi
  if [ "$i" -ge "$MIGRATE_RETRIES" ]; then
    echo "[entrypoint] migrations still failing after ${MIGRATE_RETRIES} attempts; giving up" >&2
    exit 1
  fi
  echo "[entrypoint] migrate attempt ${i} failed (DB not ready?); retrying in ${MIGRATE_RETRY_DELAY}s..."
  i=$((i + 1))
  sleep "$MIGRATE_RETRY_DELAY"
done

echo "[entrypoint] starting API..."
exec ./api
