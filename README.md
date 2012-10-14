MonkeyProxy
===========

A proxy simulates all kinds of network issues during load testing

Usage:
monkey_proxy -target www.google.com:80 -bad_rate 10 -timeout_rate 10 -dead_rate 1 --port 12345