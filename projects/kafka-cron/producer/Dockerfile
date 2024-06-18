FROM golang:1.22-bullseye as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download


COPY producer ./producer
COPY pkg ./pkg

RUN mkdir -p /out
RUN go build -o /out/producer ./producer

FROM buildpack-deps:bullseye-scm 

COPY --from=build /out/producer /producer


ENTRYPOINT ["/producer"]