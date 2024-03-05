FROM alpine:3.19.1
RUN apk --no-cache update && \ 
    apk add --no-cache busybox-extras
COPY cmd/compiled/journalingservice app
COPY static static
COPY templates templates
USER 0
CMD ["./app"]
EXPOSE 3000
# CGO_ENABLED=0 go build -o cmd/compiled/ .
# docker build -t incident:main .