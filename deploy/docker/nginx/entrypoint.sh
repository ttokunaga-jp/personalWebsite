#!/bin/sh
set -e

# PORT環境変数が未設定の場合、デフォルトで8080を使用
export PORT=${PORT:-8080}

# envsubstで設定ファイル内の${PORT}を実際の値に置換し、新しい設定ファイルを生成
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf > /etc/nginx/conf.d/generated_default.conf

# 生成した設定ファイルを使用してNginxをフォアグラウンドで実行
exec nginx -g 'daemon off;' -c /etc/nginx/conf.d/generated_default.conf
