 # === Lintas Arta's Dockerfile ===
FROM golang:alpine  AS build-env

ARG APPNAME="northstar"
ARG CONFIGPATH="/data/"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update --no-cache git gcc libc-dev tzdata;
RUN apk --no-cache add curl
#  wget gcc libc-dev make openssl py-pip;
RUN go get -u github.com/golang/dep/cmd/dep

ENV TZ=Asia/Jakarta

RUN cd $GOPATH/src/"${APPNAME}"

RUN dep ensure -v
RUN go build -v -o "${APPNAME}-res"

RUN ls -alh $GOPATH/src/
RUN ls -alh $GOPATH/src/"${APPNAME}"
RUN ls -alh $GOPATH/src/"${APPNAME}"/vendor
RUN pwd

FROM alpine

WORKDIR /go/src/
COPY --from=build-env /go/src/northstar/northstar-res /go/src/northstar
#COPY --from=build-env /go/src/northstar/deploy/conf.enc /go/src/conf.enc
COPY --from=build-env /go/src/northstar/migration/ /go/src/migration/
RUN chmod -R 775 migration
RUN ls -lrth

RUN pwd
#ENTRYPOINT /app/northstar-res
CMD ["/go/src/northstar","run"]

EXPOSE 8000
