FROM centos:latest

# Install dependencies
RUN yum -y install wget
RUN wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
RUN tar -xzf go1.10.3.linux-amd64.tar.gz
RUN mv go /usr/local
RUN yum -y install git

# Setup golang environment variables
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Install websocket dependency
RUN go get github.com/gorilla/websocket

# Add source code 
RUN mkdir /chat
COPY . /chat

# Compile source code
WORKDIR /chat/src
RUN go build main.go
CMD ./main
