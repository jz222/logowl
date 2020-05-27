FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
ENV PATH=/bin
ENV GIN_MODE=release
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD build/logowl /
CMD ["/logowl"]
EXPOSE 8080