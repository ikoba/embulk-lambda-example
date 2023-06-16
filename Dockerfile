FROM --platform=linux/amd64 amazon/aws-lambda-java:8.al2 AS java-builder
# Version of a Embulk and a Embulk Gem must be equal.
ARG embulk_version=0.11.0
COPY embulk.properties /embulk/
RUN mkdir -p /embulk/bin/
RUN yum install -y wget
RUN wget -O /embulk/bin/embulk https://dl.embulk.org/embulk-${embulk_version}.jar
RUN wget -O /embulk/bin/jruby https://repo1.maven.org/maven2/org/jruby/jruby-complete/9.3.10.0/jruby-complete-9.3.10.0.jar
# embulk, msgpack gems are required.
# liquid gem is required if you use Liquid template in a config file.
RUN java -jar /embulk/bin/embulk -X embulk_home=/embulk \
    gem install \
    embulk:${embulk_version} \
    msgpack \
    liquid \
    embulk-input-mysql \
    embulk-output-redshift -N
COPY config/* /embulk/config/

FROM golang:1.20-alpine AS go-builder
COPY src/* /go/src/
WORKDIR /go/src
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -buildid=" -trimpath -o /lambda

FROM --platform=linux/amd64 amazon/aws-lambda-java:8.al2
COPY --from=java-builder /embulk /embulk
COPY --from=go-builder /lambda /lambda

WORKDIR /
ENTRYPOINT [ "/lambda" ]
