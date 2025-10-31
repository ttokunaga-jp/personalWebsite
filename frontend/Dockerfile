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
WORKDIR /usr/share/nginx/html
COPY --from=base /app/apps/public/dist ./
COPY --from=base /app/apps/admin/dist ./admin
COPY deploy/docker/nginx/default.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
