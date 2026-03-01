#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Running smoke checks against ${BASE_URL}"

check_code() {
  local path="$1"
  local want="$2"
  local got
  got="$(curl -s -o /dev/null -w "%{http_code}" "${BASE_URL}${path}")"
  if [[ "$got" != "$want" ]]; then
    echo "FAIL: ${path} returned ${got}, expected ${want}"
    exit 1
  fi
  echo "OK: ${path} -> ${got}"
}

check_code "/" 200
check_code "/projects" 200
check_code "/writing" 200
check_code "/about" 200
check_code "/contact" 200
check_code "/projects/no-such-project" 404
check_code "/writing/no-such-post" 404
check_code "/does-not-exist" 404
check_code "/robots.txt" 200
check_code "/.well-known/security.txt" 200

health="$(curl -s "${BASE_URL}/healthz")"
if [[ "${health}" != "ok" ]]; then
  echo "FAIL: /healthz body was '${health}', expected 'ok'"
  exit 1
fi
echo "OK: /healthz body"

version_payload="$(curl -s "${BASE_URL}/version")"
if ! grep -q '"version"' <<<"${version_payload}"; then
  echo "FAIL: /version payload missing version field: ${version_payload}"
  exit 1
fi
if ! grep -q '"build_time"' <<<"${version_payload}"; then
  echo "FAIL: /version payload missing build_time field: ${version_payload}"
  exit 1
fi
echo "OK: /version payload"

if ! curl -s -D - -o /dev/null "${BASE_URL}/healthz" | grep -qi '^X-Request-Id:'; then
  echo "FAIL: missing X-Request-ID response header"
  exit 1
fi
echo "OK: X-Request-ID header present"

echo "Smoke checks passed"
