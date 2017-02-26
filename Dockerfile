FROM golang:1.7

RUN apt-get -y update && apt-get install -y fortunes

COPY . /go/src/github.com/billybobjoeaglt/matterhorn_bot/
WORKDIR /go/src/github.com/billybobjoeaglt/matterhorn_bot/

ARG version

RUN go build -ldflags "-X main.Version=$version -X main.BuildTime=`date +%s`"

EXPOSE 8080 8080
EXPOSE 8443

ENTRYPOINT ["/go/src/github.com/billybobjoeaglt/matterhorn_bot/matterhorn_bot"]
CMD ["--help"]
