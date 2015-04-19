#minimal docker file. Depends on a Postgres DB
FROM golang:latest

MAINTAINER Larry Price <larry@industrialintellect.com>

# Make sure apt is up to date
RUN apt-get update && apt-get upgrade -y
ENV LANGUAGE en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8

RUN locale-gen en_US.UTF-8 && dpkg-reconfigure locales

ENV GOROOT=/usr/src/go/

RUN CGO_ENABLED=0 go get -a -ldflags '-s' github.com/laprice/smalld/
RUN cp /gopath/src/github.com/laprice/smalld/build/Dockerfile /gopath
CMD docker build -t laprice/smalld gopath


