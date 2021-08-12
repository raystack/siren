FROM alpine:3.13
COPY siren /usr/bin/siren
EXPOSE 8080
CMD ["siren"]
