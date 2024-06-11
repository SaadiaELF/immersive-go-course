FROM golang:1.22-bullseye as build 

WORKDIR /app

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" --mount=type=bind,target=/app,source=. go mod download

RUN --mount=type=cache,target="/root/.cache/go-build" --mount=type=bind,target=/app,source=. go build -o /out/producer ./producer

FROM buildpack-deps:bullseye-scm 

COPY --from=build /out/producer /producer

EXPOSE 2112

ENTRYPOINT ["/producer"]