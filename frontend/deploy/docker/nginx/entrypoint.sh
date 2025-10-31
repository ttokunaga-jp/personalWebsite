#!/bin/sh
set -e

# PORTが未設定の場合に8080で起動させる
export PORT=${PORT:-8080}

# テンプレートから環境変数を差し替えた設定ファイルを生成してからNginxを起動する
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

exec nginx -g 'daemon off;'
