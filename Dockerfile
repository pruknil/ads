#FROM docker.artifactory.kasikornbank.com:8443/golang:1.14.6-buster as builder
FROM kbond-610-docker.artifactory.kasikornbank.com:8443/common/kbond-builder-image:2.0.1 as builder
ARG http_proxy
ARG https_proxy
ARG no_proxy

#RUN mkdir /compile
ADD ./ /usr/local/go/src/github.com/pruknil/ads
WORKDIR /usr/local/go/src/github.com/pruknil/ads
RUN go build -ldflags '-linkmode=external' -o ads
RUN file /usr/local/go/src/github.com/pruknil/ads/ads
#######################################################

#FROM docker.artifactory.kasikornbank.com:8443/golang:1.14.6-buster
FROM kbond-610-docker.artifactory.kasikornbank.com:8443/common/rhel7-mq:1.0.0
ARG http_proxy
ARG https_proxy
ARG no_proxy

#RUN apk add --no-cache ca-certificates libc6-compat file
RUN mkdir /myapp
RUN chmod 755 /myapp
COPY --from=builder /usr/local/go/src/github.com/pruknil/ads/ads /myapp/

#RUN file /myapp/ads
# Cannot add tzdata because of firewall. Copy timezone manual instead.
#RUN apk add --no-cache tzdata
COPY --from=builder /usr/local/go/src/github.com/pruknil/ads/assets/timezone/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" > /etc/timezone

WORKDIR /myapp

CMD [ "/myapp/ads" ]
