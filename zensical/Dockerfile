# syntax=docker/dockerfile:1

ARG PYTHON_IMG_TAG=3.15-slim
ARG UV_IMG_TAG=0.13.4
ARG NGINX_IMG_TAG=1.31-alpine

FROM ghcr.io/astral-sh/uv:${UV_IMG_TAG} AS uv

FROM python:${PYTHON_IMG_TAG} AS base

COPY --from=uv /uv /uvx /bin/

WORKDIR /workspace

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

FROM base AS stage

COPY zensical/pyproject.toml zensical/uv.lock zensical/

WORKDIR /workspace/zensical

RUN uv sync --locked

FROM stage AS dev

WORKDIR /workspace

ENV PATH="/workspace/zensical/.venv/bin:$PATH"

COPY zensical.toml /workspace/zensical.toml
COPY cheatsheets/ /workspace/cheatsheets/

EXPOSE 8118

CMD ["zensical", "serve", "--config-file", "/workspace/zensical.toml", "--dev-addr", "0.0.0.0:8118"]

FROM base AS build

WORKDIR /workspace

COPY zensical/pyproject.toml zensical/uv.lock zensical/

RUN uv sync \
  --locked \
  --project /workspace/zensical

COPY . .

RUN uv run \
  --project /workspace/zensical \
  zensical build \
  --config-file /workspace/zensical.toml

FROM nginxinc/nginx-unprivileged:${NGINX_IMG_TAG} AS prod

COPY --from=build /workspace/zensical/site /usr/share/nginx/html

USER 101:101

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
