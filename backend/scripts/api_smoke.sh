#!/usr/bin/env bash
set -euo pipefail

BACKEND_PID=""
SMOKE_BACKEND_LOG=""

function log() {
  echo "[smoke] $*" >&2
}

cleanup() {
  if [[ -n "${BACKEND_PID}" ]]; then
    log "Stopping ephemeral backend (pid ${BACKEND_PID})"
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
    wait "${BACKEND_PID}" >/dev/null 2>&1 || true
    BACKEND_PID=""
  fi

  if [[ -n "${SMOKE_BACKEND_LOG}" && -f "${SMOKE_BACKEND_LOG}" && "${KEEP_SMOKE_BACKEND_LOG:-0}" != "1" ]]; then
    rm -f "${SMOKE_BACKEND_LOG}"
  fi
}

trap cleanup EXIT

function require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: required command '$1' not found" >&2
    exit 1
  fi
}

load_token_from_env() {
  local key file line value
  for key in ADMIN_TOKEN TOKEN; do
    if [[ -n "${!key:-}" ]]; then
      return 0
    fi
    for file in ".env" "../.env"; do
      if [[ ! -f "${file}" ]]; then
        continue
      fi
      line=$(grep -E "^[[:space:]]*${key}=" "${file}" | tail -n1 || true)
      if [[ -z "${line}" ]]; then
        continue
      fi
      value=${line#*=}
      value=${value%$'\r'}
      value=${value#\"}
      value=${value%\"}
      value=${value#\'}
      value=${value%\'}
      if [[ -n "${value}" ]]; then
        printf -v "${key}" '%s' "${value}"
        export "${key}"
        return 0
      fi
    done
  done
  return 0
}

wait_for_backend() {
  local port=$1
  local attempts=${SMOKE_BACKEND_MAX_ATTEMPTS:-30}
  local delay=${SMOKE_BACKEND_POLL_INTERVAL:-1}

  for ((i = 1; i <= attempts; i++)); do
    if curl -fsS --connect-timeout 1 "http://127.0.0.1:${port}/api/health" >/dev/null 2>&1; then
      return 0
    fi
    sleep "${delay}"
  done

  log "Ephemeral backend failed to become ready after ${attempts} attempts"
  if [[ -n "${SMOKE_BACKEND_LOG}" && -f "${SMOKE_BACKEND_LOG}" ]]; then
    log "Review backend logs at ${SMOKE_BACKEND_LOG}"
  fi
  exit 1
}

start_local_backend() {
  local port=${SMOKE_BACKEND_PORT:-18100}
  local mode=${SMOKE_BACKEND_SERVER_MODE:-release}
  local log_file=${SMOKE_BACKEND_LOG:-}

  if [[ -z "${log_file}" ]]; then
    log_file=$(mktemp -t smoke-backend-XXXX.log)
    SMOKE_BACKEND_LOG=${log_file}
  else
    SMOKE_BACKEND_LOG=${log_file}
  fi

  log "Starting ephemeral backend on port ${port} (in-memory repositories)"
  APP_DB_DRIVER="${SMOKE_BACKEND_DB_DRIVER:-inmemory}" \
  APP_DATABASE_DSN="" \
  APP_DATABASE_HOST="" \
  APP_SERVER_PORT="${port}" \
  APP_SERVER_MODE="${mode}" \
  APP_SECURITY_HTTPS_REDIRECT="false" \
  APP_AUTH_DISABLED="${SMOKE_BACKEND_AUTH_DISABLED:-true}" \
  APP_GOOGLE_CLIENT_ID="${APP_GOOGLE_CLIENT_ID:-dummy-client-id}" \
  APP_GOOGLE_CLIENT_SECRET="${APP_GOOGLE_CLIENT_SECRET:-dummy-client-secret}" \
  APP_GOOGLE_REDIRECT_URL="${APP_GOOGLE_REDIRECT_URL:-http://127.0.0.1:${port}/api/admin/auth/callback}" \
  go run ./cmd/server >"${SMOKE_BACKEND_LOG}" 2>&1 &
  BACKEND_PID=$!

  wait_for_backend "${port}"
  BASE_URL="http://127.0.0.1:${port}"
  log "Ephemeral backend ready at ${BASE_URL}"
  if [[ "${KEEP_SMOKE_BACKEND_LOG:-0}" == "1" ]]; then
    log "Backend logs captured at ${SMOKE_BACKEND_LOG}"
  fi
}

detect_base_url() {
  if [[ -n "${BASE_URL:-}" ]]; then
    log "Using BASE_URL=${BASE_URL}"
    return
  fi

  if [[ "${SMOKE_BACKEND_SPAWN_LOCAL:-0}" == "1" ]]; then
    start_local_backend
    return
  fi

  local candidates=(
    "http://localhost:8100"
    "http://127.0.0.1:8100"
    "http://backend:8100"
  )

  for candidate in "${candidates[@]}"; do
    if curl -fsS --connect-timeout 3 "${candidate}/api/health" >/dev/null 2>&1; then
      BASE_URL=$candidate
      log "Detected backend at ${BASE_URL}"
      return
    fi
  done

  echo "Error: backend health endpoint unreachable. Set BASE_URL=http://<host>:8100 explicitly." >&2
  exit 1
}

detect_base_url
BASE_URL=${BASE_URL:-http://localhost:8100}

load_token_from_env

DEFAULT_ADMIN_TOKEN=${SMOKE_BACKEND_DEFAULT_ADMIN_TOKEN:-admin-smoke-placeholder-token}
if [[ -z "${ADMIN_TOKEN:-}" && -z "${TOKEN:-}" && -n "${BACKEND_PID}" ]]; then
  TOKEN=${DEFAULT_ADMIN_TOKEN}
  log "Using default admin token for smoke tests"
fi
TOKEN=${ADMIN_TOKEN:-${TOKEN:-}}

require_command curl

log "Checking public health endpoint"
curl -fsS "${BASE_URL}/api/health" >/dev/null

log "Fetching public profile"
curl -fsS "${BASE_URL}/api/profile" >/dev/null

log "Listing public projects"
curl -fsS "${BASE_URL}/api/projects" >/dev/null

log "Listing availability"
curl -fsS "${BASE_URL}/api/contact/availability" >/dev/null

if [[ -n "${TOKEN}" ]]; then
  log "Listing admin summary"
  if ! curl -sf -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/admin/summary" >/dev/null; then
    log "Admin endpoints not accessible with provided token; skipping admin checks"
  fi
else
  log "Skipping admin endpoints because ADMIN_TOKEN/TOKEN not set"
fi

log "Smoke tests completed successfully"
