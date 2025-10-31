#!/bin/sh
set -e

# PORT環境変数が未設定の場合、デフォルトで8080を使用
export PORT=${PORT:-8080}

# envsubstで設定ファイル内の変数を置換し、標準出力に出力
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf