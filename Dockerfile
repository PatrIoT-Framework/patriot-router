FROM ubuntu:16.04
RUN apt-get update && \
apt-get install -y --no-install-recommends locales && \
locale-gen en_US.UTF-8 && \
apt-get dist-upgrade -y && \
apt-get --purge remove openjdk* && \
 echo "oracle-java8-installer shared/accepted-oracle-license-v1-1 select true" | debconf-set-selections && \
echo "deb http://ppa.launchpad.net/webupd8team/java/ubuntu xenial main" > /etc/apt/sources.list.d/webupd8team-java-trusty.list && \
apt-key adv --keyserver keyserver.ubuntu.com --recv-keys EEA14886 && \
apt-get update && \
apt-get install -y --no-install-recommends oracle-java8-installer oracle-java8-set-default && \
apt-get clean all
RUN apt-get -y update && \
apt-get -y install git maven iptables vim curl python iproute2 python-pip iputils-ping traceroute
WORKDIR /
RUN curl -O https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz && \
tar -xvf go1.9.1.linux-amd64.tar.gz && \
mv go /usr/local
ENV PATH $PATH:/usr/local/go/bin
RUN pip install pyroute2 && \
pip install flask && \
pip install -U jsonpickle
RUN git clone https://github.com/obabec/iproute-rest.git
WORKDIR /iproute-rest
RUN git pull --force
WORKDIR /
ENV FLASK_APP /iproute-rest/app.py
RUN git clone https://github.com/obabec/iptables-api.git
RUN go get -u github.com/oxalide/go-iptables/iptables && \
go get -u github.com/abbot/go-http-auth && \
go get -u github.com/gorilla/handlers && \
go get -u github.com/gorilla/mux
WORKDIR /iptables-api/
RUN go build -o iptables-api
WORKDIR /
ENTRYPOINT flask run --host=0.0.0.0
#ENTRYPOINT tail -f /dev/null
WORKDIR /iptables-api
