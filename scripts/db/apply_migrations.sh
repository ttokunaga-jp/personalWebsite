#!/usr/bin/env bash
set -euo pipefail

# Verifies that deploy/mysql/schema.sql and migrations apply cleanly against a MySQL 8.0 instance.
# Used in CI but can also be executed locally.

MYSQL_IMAGE="${MYSQL_IMAGE:-mysql:8.0}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-root}"
DATABASE_NAME="${DATABASE_NAME:-app}"
CONTAINER_NAME="schema-check-$(date +%s)"

cleanup() {
  docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true
}

trap cleanup EXIT

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required to run this script." >&2
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "docker daemon is not reachable. Start Docker Desktop or ensure the current user can access /var/run/docker.sock." >&2
  exit 1
fi

docker run \
  --rm \
  --name "${CONTAINER_NAME}" \
  -e MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD}" \
  -e MYSQL_DATABASE="${DATABASE_NAME}" \
  -d \
  "${MYSQL_IMAGE}" \
  --default-authentication-plugin=mysql_native_password \
  --log_bin_trust_function_creators=1 \
  --character-set-server=utf8mb4 \
  --collation-server=utf8mb4_unicode_ci \
  >/dev/null

echo "Waiting for MySQL container ${CONTAINER_NAME} to become ready..."
for attempt in {1..30}; do
  if docker exec "${CONTAINER_NAME}" mysqladmin ping -uroot -p"${MYSQL_ROOT_PASSWORD}" --silent; then
    break
  fi
  sleep 2
  if [[ ${attempt} -eq 30 ]]; then
    echo "MySQL container failed to start" >&2
    exit 1
  fi
done

apply_sql() {
  local file="$1"
  echo "Applying ${file}"
  docker exec -i "${CONTAINER_NAME}" mysql -uroot -p"${MYSQL_ROOT_PASSWORD}" "${DATABASE_NAME}" < "${file}"
}

apply_sql "deploy/mysql/schema.sql"

if [[ -d "deploy/mysql/migrations" ]]; then
  while IFS= read -r -d '' file; do
    apply_sql "${file}"
  done < <(find deploy/mysql/migrations -maxdepth 1 -type f -name "*.sql" -print0 | sort -z)
fi

echo "Schema and migrations applied successfully."
