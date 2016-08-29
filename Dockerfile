FROM golang:1.6

COPY . /go/src/github.com/billybobjoeaglt/matterhorn_bot/
WORKDIR /go/src/github.com/billybobjoeaglt/matterhorn_bot/
RUN make build

EXPOSE 8080 8080

ENTRYPOINT ["/bin/bash", "-c"]
CMD ["/bin/bash"]