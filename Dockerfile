FROM golang:alpine

ARG APPNAME="northstar"
ARG CONFIGPATH="$$GOPATH/src/northstar"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev tzdata;
#  tzdata wget gcc libc-dev make openssl py-pip;

ENV TZ=Asia/Jakarta

RUN go get -u github.com/golang/dep/cmd/dep

CMD if [ "${APPENV}" = "staging" ] || [ "${APPENV}" = "production" ] ; then \
        openssl aes-256-cbc -d -in deploy/conf.enc -out config.yaml -pbkdf2 -pass file:./public.pem ; \
    elif [ "${APPENV}" = "dev" ] ; then \
        cp deploy/dev-config.yaml config.yaml ; \
    fi \
    && dep ensure -v \
    && go build -v -o $GOPATH/bin/"${APPNAME}" \
    && "${APPNAME}" run \
EXPOSE 8000