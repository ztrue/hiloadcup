FROM ztrue/hlc

COPY . /go/src/app
RUN mkdir -p /tmp/unzip
COPY ./data.zip /tmp/data/data.zip

#RUN ffjson /go/src/app/structs/structs.go

RUN go install

RUN chmod +x /go/src/app/run.sh

CMD bash /go/src/app/run.sh
