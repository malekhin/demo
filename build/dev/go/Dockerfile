FROM golang:1.23-bullseye

ARG GIT_USERNAME
ARG GIT_PASSWORD

ENV XDG_RUNTIME_DIR=/tmp
ENV GOPRIVATE=gitlab.inguru.dev

RUN apt update -y && \
    DEBIAN_FRONTEND=noninteractive apt install -y git openssl build-essential xorg libssl-dev wkhtmltopdf pgloader pip

RUN go install github.com/cosmtrek/air@v1.40.4
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@v1.24.2
RUN go install -tags 'postgres,clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.2
RUN go install github.com/vektra/mockery/v2@v2.40.2
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.61.0
RUN go install github.com/boumenot/gocover-cobertura@v1.3.0
RUN pip install diff_cover==9.2

WORKDIR /app

COPY ./migrate_pg.sh /migrate_pg.sh
COPY ./entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
