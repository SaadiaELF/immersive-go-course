FROM golang:1.22-bullseye as build 

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download


COPY consumer ./consumer
COPY pkg ./pkg

RUN mkdir -p /out
RUN go build -o /out/consumer ./consumer

FROM buildpack-deps:bullseye-scm 

COPY --from=build /out/consumer /consumer


ENTRYPOINT ["/consumer"]