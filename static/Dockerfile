# Download API server
FROM alpine:latest as linux

RUN apk add git &&\
    git clone https://github.com/g3un/sp1-backend.git /root/src/

# Build API server
FROM golang:latest as build

COPY --from=linux /root/src/ /root/src/
RUN cd /root/src/ &&\
    go build -o /root/bin/main /root/src/main.go

# Start API server
FROM chromedp/headless-shell:latest

COPY --from=build /root/bin/main /root/main
RUN apt-get update &&\
    apt-get upgrade -y &&\
    apt install dumb-init -y

ENTRYPOINT ["dumb-init", "--"]
CMD ["/root/main"]
