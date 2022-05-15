FROM scratch
ENTRYPOINT ["/demo-webhook-app"]
COPY demo-webhook-app /
COPY config.yml /