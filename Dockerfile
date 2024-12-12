FROM --platform=$BUILDPLATFORM golang:1.23 AS build

WORKDIR /src
COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o woodpecker-jsonnet cmd/woodpecker-jsonnet/main.go

FROM --platform=$BUILDPLATFORM debian:bookworm-slim
ENV GODEBUG=netdns=go

# copy certs from build image
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# copy agent binary
COPY --from=build /src/woodpecker-jsonnet /bin/

ENTRYPOINT ["/bin/woodpecker-jsonnet"]
