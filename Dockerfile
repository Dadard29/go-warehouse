FROM golang:1.13-alpine

ARG ARG_VERSION
ARG ARG_HOST_SUB
ARG DB_USER
ARG DB_PASSWORD

ENV VERSION=$ARG_VERSION
ENV HOST_SUB=$ARG_HOST_SUB
ENV DB_USER=$ARG_DB_USER
ENV DB_PASSWORD=$ARG_DB_PASSWORD

ENV CORS_ORIGIN=https://dadard.fr

RUN apk add --update git gcc libc-dev

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]