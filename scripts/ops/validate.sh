#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TF_ENV_DIR="${ROOT_DIR}/terraform/environments/dev"

echo "==> terraform fmt check"
terraform fmt -recursive -check "${ROOT_DIR}"

echo "==> terraform validate (dev environment)"
terraform -chdir="${TF_ENV_DIR}" init -backend=false >/dev/null
terraform -chdir="${TF_ENV_DIR}" validate

if command -v tflint >/dev/null 2>&1; then
  echo "==> tflint"
  (cd "${TF_ENV_DIR}" && tflint --init >/dev/null && tflint)
else
  echo "==> skipping tflint (tflint not found in PATH)"
fi

if command -v yamllint >/dev/null 2>&1; then
  echo "==> yamllint (ops YAML assets)"
  yamllint "${ROOT_DIR}/deploy" "${ROOT_DIR}/.github/workflows"
else
  echo "==> skipping yamllint (yamllint not found in PATH)"
fi

if command -v gcloud >/dev/null 2>&1 && [[ -n "${OPS_GCP_PROJECT:-}" ]]; then
  echo "==> gcloud monitoring policies lint"
  gcloud beta monitoring policies lint --project "${OPS_GCP_PROJECT}" >/dev/null
else
  echo "==> skipping gcloud monitoring lint (set OPS_GCP_PROJECT and ensure gcloud is installed)"
fi

echo "ops validation completed successfully."

