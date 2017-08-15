FROM ztrue/hlc

COPY . /go/src/app
COPY ./data.zip /tmp/data/data.zip

#RUN ffjson /go/src/app/structs/structs.go
