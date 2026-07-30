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
  # Capture the output so a failure can be classified instead of blamed on the
  # database. `if cmd` disables set -e for cmd, so a failing migrate does not
  # exit here.
  if migrate_out="$(./migrate up 2>&1)"; then
    printf '%s\n' "$migrate_out"
    break
  fi
  printf '%s\n' "$migrate_out" >&2

  # A configuration error is not a readiness race: retrying it 30 times only
  # buries the real message under a misleading "DB not ready?" loop while
  # `restart: unless-stopped` keeps the container alive and unhealthy. Fail fast
  # so the actual cause is the last thing in the log.
  case "$migrate_out" in
    *"config:"*|*"must be set"*)
      echo "[entrypoint] migrate failed on configuration, not database readiness — not retrying" >&2
      exit 1
      ;;
  esac

  if [ "$i" -ge "$MIGRATE_RETRIES" ]; then
    echo "[entrypoint] migrations still failing after ${MIGRATE_RETRIES} attempts; giving up" >&2
    exit 1
  fi
  echo "[entrypoint] migrate attempt ${i} failed (database may not be ready); retrying in ${MIGRATE_RETRY_DELAY}s..."
  i=$((i + 1))
  sleep "$MIGRATE_RETRY_DELAY"
done

echo "[entrypoint] starting API..."
exec ./api
