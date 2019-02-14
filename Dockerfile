FROM golang:1.11

RUN go get -d -v github.com/olivere/elastic && \
go get -d -v github.com/sirupsen/logrus && \
go get -d -v gopkg.in/sohlich/elogrus.v3 &&\
go get -d -v github.com/abbot/go-http-auth && \
go get -d -v github.com/gorilla/handlers && \
go get -d -v github.com/gorilla/mux &&\
go get -d -v github.com/oxalide/go-iptables/iptables

WORKDIR /go/src
ADD api/iprouteRESt /go/src/api/iprouteRESt
ADD api/iptables-api /go/src/api/iptables-api
ADD api/iproute2 /go/src/api/iproute2
WORKDIR /

WORKDIR /go/src/api/iptables-api/
RUN go build  -o /api/iptables-api main.go 
WORKDIR /go/src/api/iprouteRESt/
RUN go build -o /api/iprouteRESt Controller.go 

ADD start_daemons /start_daemons

WORKDIR /
ENTRYPOINT ["/start_daemons"]
CMD ["/api/iptables-api", "-ip", "0.0.0.0"]
