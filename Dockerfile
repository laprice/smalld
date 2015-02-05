#minimal docker file. Depends on a Postgres DB
FROM golang:latest

MAINTAINER Larry Price <larry@industrialintellect.com>

# Make sure apt is up to date
RUN apt-get update && apt-get upgrade -y && apt-get install -y language-pack-en
ENV LANGUAGE en_US.UTF-8
ENV LANG en_US.UTF-8
ENV LC_ALL en_US.UTF-8

RUN locale-gen en_US.UTF-8 && dpkg-reconfigure locales

ENV GOROOT=/usr/src/go/

RUN git clone https://github.com/laprice/smalld.git $GOROOT/src/github.com/laprice/ \
     && cd $GOROOT/src/github.com/laprice/smalld && make 

