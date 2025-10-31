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
COPY deploy/docker/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY deploy/docker/nginx/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN mkdir -p /var/cache/nginx /var/run/nginx /var/log/nginx \
    && addgroup -g 65532 -S cloudrun \
    && adduser -S -D -H -u 65532 -s /sbin/nologin cloudrun -G cloudrun \
    && chown -R cloudrun:cloudrun /usr/share/nginx/html /etc/nginx /var/cache/nginx /var/run/nginx /var/log/nginx

USER 65532:65532
EXPOSE 8080

CMD ["/entrypoint.sh"]
