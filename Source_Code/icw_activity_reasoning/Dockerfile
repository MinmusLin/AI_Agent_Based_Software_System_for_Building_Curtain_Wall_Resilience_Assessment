FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS go-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY --from=icw_common go.mod go.sum ./icw_common/
COPY go.mod go.sum ./icw_activity_reasoning/

WORKDIR /src/icw_activity_reasoning

RUN --mount=type=cache,target=/go/pkg/mod apk add --no-cache ca-certificates && go mod download

WORKDIR /src

COPY --from=icw_common . ./icw_common
COPY cmd ./icw_activity_reasoning/cmd
COPY configs ./icw_activity_reasoning/configs
COPY internal ./icw_activity_reasoning/internal
COPY rpc ./icw_activity_reasoning/rpc

WORKDIR /src/icw_activity_reasoning

RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	go build -trimpath -buildvcs=false -tags=netgo,osusergo -ldflags="-s -w -buildid=" -o /out/icw_activity_reasoning ./cmd/server

FROM python:3.13.13-slim AS runtime

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV UV_CACHE_DIR=/tmp/uv-cache
ENV UV_PROJECT_ENVIRONMENT=/app/.venv
ENV REASONING_WORK_DIR=/app/.runtime

WORKDIR /app

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates curl libgl1 libglib2.0-0 libgomp1 \
	&& rm -rf /var/lib/apt/lists/* \
	&& curl -LsSf https://astral.sh/uv/install.sh | sh \
	&& mv /root/.local/bin/uv /usr/local/bin/uv \
	&& chmod +x /usr/local/bin/uv

COPY pyproject.toml uv.lock ./

RUN --mount=type=cache,target=/tmp/uv-cache UV_LINK_MODE=copy uv sync --frozen --no-dev --no-install-project

COPY python ./python
COPY --from=go-builder /out/icw_activity_reasoning /app/icw_activity_reasoning

RUN mkdir -p /app/.runtime /tmp/stain-detection-cache

EXPOSE 8003

ENTRYPOINT ["/app/icw_activity_reasoning"]
