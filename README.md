MonkeyProxy
===========

A proxy simulates all kinds of network issues during load testing

##Installation
* download go installer from: http://code.google.com/p/go/downloads/list
* compile: cd MonkeyProxy; export GO_PATH=\`pwd\`; go build src/monkey_proxy

##Prerequisite
on test application server, update /etc/hosts to point backend servers to monkey proxy server

##Usage
launch monkey_proxy with the following command:
`monkey_proxy -target www.google.com:80  -port 12345 -bad_rate 10 -timeout_rate 10 -dead_rate 1`
- -target: the target server hostname and port *required*
- -port: local port, typically it should be same as target server port
- -bad_rate: percentage of requests will return partial response
- -timeout_rate: percentage of requests will be timed out
- -timeout: the timed out requests will wait exact *timeout* milliseonds before response