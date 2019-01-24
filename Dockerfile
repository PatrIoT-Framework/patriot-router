FROM ubuntu:16.04
RUN apt-get update && \
apt-get install -y --no-install-recommends locales && \
locale-gen en_US.UTF-8 && \
apt-get dist-upgrade -y
RUN apt-get -y update && \
apt-get -y install git maven iptables vim curl python iproute2 python-pip iputils-ping traceroute
WORKDIR /
RUN curl -O https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz && \
tar -xvf go1.9.1.linux-amd64.tar.gz && \
mv go /usr/local
ENV PATH $PATH:/usr/local/go/bin

ADD api/iprouteRESt /api/iprouteRESt
ADD api/iptables-api /api/iptables-api
ADD api/iproute2 /usr/local/go/src/iproute2
ADD start_daemons /start_daemons
RUN go get -u github.com/olivere/elastic && \
go get -u github.com/sirupsen/logrus && \
go get -u gopkg.in/sohlich/elogrus.v3 &&\
go get -u github.com/abbot/go-http-auth && \
go get -u github.com/gorilla/handlers && \
go get -u github.com/gorilla/mux &&\
go get -u github.com/oxalide/go-iptables/iptables

WORKDIR /api/iptables-api/
RUN go build -o iptables-api

WORKDIR /api/iprouteRESt/
RUN go build -o Controller.go
WORKDIR /
CMD ["tail","-f","/dev/null"]
ENTRYPOINT /start_daemons
