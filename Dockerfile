FROM centos:latest

# Install dependencies
RUN yum -y install wget
RUN wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz
RUN tar -xzf go1.10.3.linux-amd64.tar.gz
RUN mv go /usr/local

# Setup golang environment variables
ENV GOROOT=/usr/local/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Add source code 
RUN install -d /usr/local/go/src/LOG735-PG/src
COPY src/ /usr/local/go/src/LOG735-PG/src

# Compile source code
WORKDIR /usr/local/go/src/LOG735-PG/src
RUN go build main.go

CMD /usr/local/go/src/LOG735-PG/src/main
