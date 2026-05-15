# Stage 1: build a fully static Linux binary.
# CGO is disabled so the resulting executable has zero shared-library
# dependencies and runs on the distroless static base in stage 2.
FROM golang:1.26-alpine AS build

WORKDIR /app

# Cache module downloads as a separate layer so iteration on source code
# doesn't re-pull the dependency graph.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -ldflags="-s -w" strips the symbol table and DWARF debug info — saves
# ~30% binary size with no runtime impact (panics still print a stack;
# only `go tool addr2line` style introspection is affected).
# -trimpath removes the build host's path prefix from embedded source
# locations, making builds reproducible.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /out/velosearch \
    ./cmd/server

# Stage 2: minimal runtime.
# distroless/static-debian12:nonroot ships only ca-certificates, /etc/passwd,
# tzdata, and a non-root user (uid 65532). No shell, no package manager,
# no busybox — vastly reduced attack surface.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/velosearch /velosearch

# Default data directory; bind-mount a volume here for persistent WAL.
WORKDIR /data

EXPOSE 50051

USER nonroot

ENTRYPOINT ["/velosearch"]
CMD ["-addr=:50051", "-data-dir=/data"]
