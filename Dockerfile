FROM golang:1.10.3-stretch

ENV DBHOST="pg"
ENV DBPORT="5432"
ENV DBUSER="postgres"
ENV DBPASS="secretpassword"
ENV DBNAME="api_development"

ENV APP_NAME="api-go"

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY . /usr/src/app/

EXPOSE 1323

ENTRYPOINT ["./Server"]
