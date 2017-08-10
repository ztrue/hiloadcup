FROM golang:onbuild

COPY ./data.zip /tmp/data/data.zip

EXPOSE 80
