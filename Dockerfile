FROM alpine:3.15
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/demo-webhook-app"]
COPY demo-webhook-app /
COPY config.yml /
EXPOSE 8000