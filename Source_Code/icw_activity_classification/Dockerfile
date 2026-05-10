FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS go-builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY --from=icw_common go.mod go.sum ./icw_common/
COPY go.mod go.sum ./icw_activity_classification/

WORKDIR /src/icw_activity_classification

RUN --mount=type=cache,target=/go/pkg/mod apk add --no-cache ca-certificates && go mod download

WORKDIR /src

COPY --from=icw_common . ./icw_common
COPY . ./icw_activity_classification

WORKDIR /src/icw_activity_classification

RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	go build -trimpath -buildvcs=false -tags=netgo,osusergo -ldflags="-s -w -buildid=" -o /out/icw_activity_classification ./cmd/server

FROM scratch

WORKDIR /app

COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=go-builder /out/icw_activity_classification /app/icw_activity_classification

EXPOSE 8002

USER 65532:65532

ENTRYPOINT ["/app/icw_activity_classification"]
