FROM golang:1.11

WORKDIR /go/src/app
ADD api/iprouteRESt /go/src/app/api/iprouteRESt
ADD api/iptables-api /go/src/app/api/iptables-api
ADD api/iproute2 /usr/local/go/src/iproute2
WORKDIR /
ADD start_daemons /start_daemons

RUN go get -d -v github.com/olivere/elastic && \
go get -d -v github.com/sirupsen/logrus && \
go get -d -v gopkg.in/sohlich/elogrus.v3 &&\
go get -d -v github.com/abbot/go-http-auth && \
go get -d -v github.com/gorilla/handlers && \
go get -d -v github.com/gorilla/mux &&\
go get -d -v github.com/oxalide/go-iptables/iptables

WORKDIR /go/src/app/api/iptables-api/
RUN go build -o main.go
WORKDIR /go/src/app/api/iprouteRESt/
RUN go build -o Controller.go
WORKDIR /
CMD ["tail","-f","/dev/null"]
ENTRYPOINT /start_daemons
