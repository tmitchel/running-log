FROM golang

ADD . /Users/tmitchel/Documents/projects/running-log

RUN go build -o running-log -v ./cmd/running-log/main.go

ENTRYPOINT [ "./running-log" ]

EXPOSE 8000