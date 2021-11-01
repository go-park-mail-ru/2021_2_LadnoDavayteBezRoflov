# syntax=docker/dockerfile:1
##############################
# Dockerfile to run backend  #
# Based on golang:latest     #
##############################

FROM golang:latest

WORKDIR /backend

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest && go mod tidy

COPY . .

CMD /bin/bash -c '$GOPATH/bin/CompileDaemon -log-prefix=false -polling=true -polling-interval=500 -build="go build ./cmd/api/" -command="./api"'

EXPOSE 8000

LABEL \
      name="2021_2_LadnoDavayteBezRoflov_Backend" \
      description="Launch LadnoDavayteBezRoflov_Backend" \
      version="1.0" \
      maintainer=""
