FROM postgres:13-alpine

RUN apk add --update htop
#RUN mkdir -p /data
#COPY ./data/* /data/
COPY ./init.sql /docker-entrypoint-initdb.d/
#RUN chown postgres /data
#RUN chmod +w /data