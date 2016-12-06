FROM golang:1.6

ARG version

RUN apt-get -y update && apt-get install -y fortunes

COPY . /go/src/github.com/billybobjoeaglt/matterhorn_bot/
WORKDIR /go/src/github.com/billybobjoeaglt/matterhorn_bot/
RUN VERSION=$version make build

EXPOSE 8080 8080

ENTRYPOINT ["/bin/bash", "-c"]
CMD ["/bin/bash"]
