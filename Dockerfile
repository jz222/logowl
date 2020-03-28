FROM scratch
ENV GIN_MODE=release
ADD build/loggy /
CMD ["/loggy"]
EXPOSE 8080