# Context dir .
FROM golang:1.15 as build
#ENV GOPROXY
ENV GOSUMDB off
ENV GO111MODULE on

ENV MG_PORT 8080
ENV MG_DB_HOST 0.0.0.0
ENV MG_DB_PORT 5432
ENV MG_DB_USER postgres
ENV MG_DB_PASS xxx
ENV MG_DB_KEY xxx
ENV MG_DB_NAME postgres

WORKDIR /go/src/magnus
COPY /src /go/src/magnus
RUN go mod download
WORKDIR /go/src/magnus/cmd/mgquiz
RUN GOOS=linux CGO_ENABLED=0 go build -o magnus

ENTRYPOINT ["./magnus"]
