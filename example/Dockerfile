FROM golang:alpine AS base

# DEVELOPMENT
FROM base AS development
WORKDIR /tycho

COPY . ./
COPY /example ./example

WORKDIR /tycho/example

RUN go mod download
RUN go install github.com/cespare/reflex@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

RUN go install github.com/wearepointers/sqlboiler-erg@latest
RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

COPY /example ./example

CMD ["reflex","-s", "--", "go", "run", "cmd/server/main.go"]
