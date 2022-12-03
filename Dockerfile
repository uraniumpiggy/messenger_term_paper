FROM scratch
COPY server/config.yml /
COPY server/cmd/main/main /
CMD ["/main"]