FROM node:20-alpine AS base
ENV PNPM_HOME=/pnpm
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app

COPY package.json pnpm-workspace.yaml tsconfig.json tsconfig.base.json ./
COPY pnpm-lock.yaml ./pnpm-lock.yaml
COPY packages ./packages
COPY apps ./apps

RUN pnpm install
RUN pnpm -r build

FROM nginx:1.25-alpine
RUN apk add --no-cache gettext
WORKDIR /usr/share/nginx/html
COPY --from=base /app/apps/public/dist ./
COPY --from=base /app/apps/admin/dist ./admin
RUN mkdir -p /etc/nginx/templates
COPY deploy/docker/nginx/default.conf.template /etc/nginx/templates/default.conf.template

ENV NGINX_ENVSUBST_TEMPLATE_DIR=/etc/nginx/templates \
    NGINX_ENVSUBST_OUTPUT_DIR=/etc/nginx/conf.d \
    PORT=8080 \
    API_PROXY_PASS=http://backend:8100/api/

RUN mkdir -p \
        /var/cache/nginx/client_temp \
        /var/cache/nginx/proxy_temp \
        /var/cache/nginx/fastcgi_temp \
        /var/cache/nginx/uwsgi_temp \
        /var/cache/nginx/scgi_temp \
        /var/run/nginx \
        /var/log/nginx \
        /var/lib/nginx/tmp/client_body \
        /var/lib/nginx/tmp/proxy \
        /var/lib/nginx/tmp/fastcgi \
        /var/lib/nginx/tmp/uwsgi \
        /var/lib/nginx/tmp/scgi \
    && chown -R nginx:nginx /usr/share/nginx/html /etc/nginx /var/cache/nginx /var/run/nginx /var/log/nginx /var/lib/nginx
EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
