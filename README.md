A distributed scrabble cheating tool written in Go. It really has no reason to exist.
I just wanted to try out go, zeromq and protocol buffers respectively.

## Dependencies
    * zeromq 4 http://zeromq.org/intro:get-the-software
    * zmq4 https://github.com/pebbe/zmq4
    * protocol buffers https://github.com/google/protobuf
    * protobuf https://github.com/golang/protobuf

## Setting up curvezmq
Both the client and server generate their keys in the current working
directory unless specified otherwise.

To authorize a client to access the server copy its public key to a
known directory and point the server to that directory.
