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
# Add functional tests
RUN install -d /usr/local/go/src/LOG735-PG/ft
COPY ft/ /usr/local/go/src/LOG735-PG/ft
# Start testing
#WORKDIR /usr/local/go/src/LOG735-PG/src
# Test client package
#RUN pushd client; go test; popd
# Test miner package
#RUN pushd miner; go test; popd
# Run functional tests
#RUN install -d /usr/local/go/src/LOG735-PG/scripts
#COPY scripts/iterate_dir.sh /usr/local/go/src/LOG735-PG/scripts
#RUN /usr/local/go/src/LOG735-PG/scripts/iterate_dir.sh
# Build image
WORKDIR /usr/local/go/src/LOG735-PG/src
RUN go build main.go

CMD /usr/local/go/src/LOG735-PG/src/main
