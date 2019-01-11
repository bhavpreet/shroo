# Dockerfile for cross compiling
FROM mitchallen/pi-cross-compile
RUN apt-get update
RUN apt-get install wget
RUN apt-get install file
RUN cd /tmp/
RUN wget https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.11.4.linux-amd64.tar.gz
RUN echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
RUN mkdir /go
RUN echo "export GOPATH=/go" >> ~/.bashrc
WORKDIR /go/src/
    


