FROM postgres:13-alpine

RUN apk add --update htop
COPY ./init.sql /docker-entrypoint-initdb.d/