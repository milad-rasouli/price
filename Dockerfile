# sudo docker build -t price:0.0.1 .
# sudo docker run --rm -p 8080:8080 --name price price:0.0.1
FROM golang:1.25.1 AS base
COPY . /app

FROM base AS build
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build ./cmd/.

FROM debian:bookworm-slim AS prod
WORKDIR /app
COPY --from=build /app/build /app/build
RUN chmod a+x /app/build
ENTRYPOINT ["/app/build"]
