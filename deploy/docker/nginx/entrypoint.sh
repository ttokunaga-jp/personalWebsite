#!/bin/sh
set -ex

export PORT=${PORT:-8080}
echo "--- Substituting PORT: ${PORT} ---"

echo "--- Checking permissions ---"
ls -ld /etc/nginx/conf.d
ls -l /etc/nginx/conf.d/default.conf.template

echo "--- Generating nginx config ---"
envsubst '${PORT}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

echo "--- Generated /etc/nginx/conf.d/default.conf ---"
cat /etc/nginx/conf.d/default.conf
echo "---------------------------------------------"

echo "--- Starting nginx ---"
nginx -g 'daemon off;'
echo "--- nginx exited with code $? ---"