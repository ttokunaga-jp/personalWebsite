#!/usr/bin/env bash
set -euo pipefail

function log() {
  echo "[smoke] $*" >&2
}

function require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Error: required command '$1' not found" >&2
    exit 1
  fi
}

detect_base_url() {
  if [[ -n "${BASE_URL:-}" ]]; then
    log "Using BASE_URL=${BASE_URL}"
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
  curl -fsS -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/admin/summary" >/dev/null
else
  log "Skipping admin endpoints because ADMIN_TOKEN/TOKEN not set"
fi

log "Smoke tests completed successfully"
