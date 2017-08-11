FROM golang:onbuild

#RUN apt-get install software-properties-common
#RUN add-apt-repository ppa:chris-lea/redis-server
#RUN apt-get install redis-server
#RUN redis-server --version

COPY ./data.zip /tmp/data/data.zip

EXPOSE 80
