FROM alpine:3.13
COPY siren .
EXPOSE 8080
ENTRYPOINT ["./siren"]