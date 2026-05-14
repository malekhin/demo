FROM golang:1.22.4-alpine-base1 as build

WORKDIR /app
COPY . ./
RUN go build -o ./tmp/app ./src/cmd/app

FROM registry.inguru.dev/base_images/golang:1.22.4-alpine-work1

WORKDIR /app

ENV LOGGER_LEVEL=debug


ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 8000

COPY ./config.dev.yaml.dist ./config.dev.yaml
COPY config.prod.yaml ./
COPY docs ./docs
COPY --from=build /app/tmp/app /app

CMD ./app