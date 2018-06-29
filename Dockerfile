FROM centos:latest

# Install dependencies
RUN yum -y install wget
RUN wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
RUN tar -xzf go1.10.3.linux-amd64.tar.gz
RUN mv go /usr/local

# Setup golang environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/root/LOG735
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Add source code 
RUN install -d /root/LOG735
COPY src/ /root/LOG735

# Compile source code
WORKDIR /root/LOG735
RUN go build main.go

CMD /root/LOG735/main
