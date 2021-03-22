FROM golang:1.15
WORKDIR /go/src/github.com/odpf/siren
COPY . .
RUN make dist
EXPOSE 3000
ENTRYPOINT ["/go/src/github.com/odpf/siren/dist/linux-amd64/siren"]
