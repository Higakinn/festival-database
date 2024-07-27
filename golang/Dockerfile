## build stage
FROM golang:1.20.0-bullseye as deploy-builder

WORKDIR /app

# module file copy
COPY src/go.mod src/go.sum ./
RUN go mod download

# application code copy
COPY ./src . 

RUN go build -trimpath -ldflags "-w -s" -o app

## deploy stage
FROM debian:bullseye-slim as deploy

RUN apt-get update && apt-get install -y ca-certificates openssl

COPY --from=deploy-builder /app/app .

CMD ["./app"]

## dev stage
FROM golang:1.21.0 as dev
WORKDIR /app

RUN go install github.com/air-verse/air@latest && \
  go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
  go install github.com/ramya-rao-a/go-outline@latest && \
  go install github.com/nsf/gocode@latest && \
  go install github.com/acroca/go-symbols@latest && \
  go install github.com/fatih/gomodifytags@latest && \
  go install github.com/josharian/impl@latest && \
  go install github.com/haya14busa/goplay/cmd/goplay@latest && \
  go install github.com/go-delve/delve/cmd/dlv@latest && \
  go install golang.org/x/lint/golint@latest && \
  go install golang.org/x/tools/gopls@latest

CMD ["air", "-c", "air.toml"]

FROM golang:1.21.0 as ci
WORKDIR /app

CMD ["tail", "-f", "/dev/null"]
