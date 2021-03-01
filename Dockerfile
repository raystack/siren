FROM golang:1.15
WORKDIR /go/src/github.com/odpf/siren
COPY . .
RUN make test
RUN make dist

FROM alpine:latest
RUN ["apk", "update"]
RUN ["apk", "add", "libc6-compat"]
WORKDIR /root/
COPY --from=0 /go/src/github.com/odpf/siren .
EXPOSE 3000
ENTRYPOINT ["/root/dist/linux-amd64/siren"]
