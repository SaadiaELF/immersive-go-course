FROM golang:1.22-bullseye as build 

WORKDIR /app

ENV GOCACHE=/root/.cache/go-build

RUN --mount=type=cache,target="/root/.cache/go-build" --mount=type=bind,target=/app,source=. go mod download

RUN --mount=type=cache,target="/root/.cache/go-build" --mount=type=bind,target=/app,source=. go build -o /out/consumer ./consumer

FROM buildpack-deps:bullseye-scm 

COPY --from=build /out/consumer /consumer

ENTRYPOINT ["/consumer"]