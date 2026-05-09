FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY --from=icw_common go.mod go.sum ./icw_common/
COPY go.mod go.sum ./icw_activity_summary/

WORKDIR /src/icw_activity_summary

RUN --mount=type=cache,target=/go/pkg/mod apk add --no-cache ca-certificates && go mod download

WORKDIR /src

COPY --from=icw_common . ./icw_common
COPY . ./icw_activity_summary

WORKDIR /src/icw_activity_summary

RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	go build -trimpath -buildvcs=false -tags=netgo,osusergo -ldflags="-s -w -buildid=" -o /out/icw_activity_summary ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/icw_activity_summary /app/icw_activity_summary

EXPOSE 8000

USER 65532:65532

ENTRYPOINT ["/app/icw_activity_summary"]
