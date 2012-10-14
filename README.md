MonkeyProxy
===========

A proxy simulates all kinds of network issues during load testing

#Installation
* download go installer from: http://code.google.com/p/go/downloads/list
* compile: cd MonkeyProxy; export GO_PATH=`pwd`;go build src/http_proxy -o monkey_proxy

#Usage
monkey_proxy -target www.google.com:80 -bad_rate 10 -timeout_rate 10 -dead_rate 1 --port 12345