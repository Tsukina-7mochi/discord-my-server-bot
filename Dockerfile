FROM ubuntu:24.04 AS base

RUN apt update && apt upgrade -y
RUN apt install ca-certificates -y
RUN update-ca-certificates


FROM base AS dev

RUN apt install -y curl wget tar

WORKDIR /tmp
RUN wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
ENV GOPATH /root/go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

WORKDIR /app
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
    go install github.com/air-verse/air@latest
RUN --mount=type=cache,target=/root/go/pkg/mod,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

CMD ["/root/go/bin/air"]


FROM dev AS build

RUN --mount=type=bind,source=.,target=. \
    go build -o /app-build/main ./cmd/bot.go


FROM base AS prod

WORKDIR /app
COPY --from=build /app-build/main /app/main
CMD ["/bin/bash", "-c", "/app/main"]
