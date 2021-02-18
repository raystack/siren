FROM alpine:latest
RUN ["apk", "update"]
RUN ["apk", "add", "libc6-compat"]

WORKDIR /opt/
COPY ./dist/siren/linux-amd64/siren .
RUN chmod +x ./siren

EXPOSE 3000
ENTRYPOINT ["/opt/siren"]
