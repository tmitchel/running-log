FROM golang:1.13

WORKDIR /go/src/running-log
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD [ "running-log" ]

EXPOSE 8000
EXPOSE 6379