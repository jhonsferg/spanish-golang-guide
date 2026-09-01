# Dockerfile de referencia para el Capítulo 46: build multi-stage,
# imagen final mínima (distroless/scratch), sin toolchain de Go en runtime.

FROM golang:1.27 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/servidor .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/servidor /servidor
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/servidor"]
