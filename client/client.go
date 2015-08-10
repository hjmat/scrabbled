/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

/*
 * scrabbled client
 * USAGE: client [--port <PORT>] [--host <HOSTNAME>]
 */

package main

import (
	"github.com/hjmat/scrabbled/condlog"
	scrabble "github.com/hjmat/scrabbled/proto"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"bufio"
	"flag"
	"fmt"
	"os"
)

func solve(hand string, sock *zmq.Socket) []string {
	req := &scrabble.Request{Hand: hand}

	reqMsg, err := proto.Marshal(req)
	condlog.Fatal(err, "Unable to marshal message")

	_, err = sock.Send(string(reqMsg), 0)
	condlog.Fatal(err, "Unable to marshal message")

	resMsg, err := sock.Recv(0)
	condlog.Fatal(err, "Unable to receive response")

	res := &scrabble.Response{}
	err = proto.Unmarshal([]byte(resMsg), res)
	condlog.Fatal(err, "Unable to unmarshal request")

	return res.Words
}

func main() {
	hostPtr := flag.String("host", "localhost", "host")
	portPtr := flag.Int("port", 30000, "port")
	flag.Parse()

	sock, err := zmq.NewSocket(zmq.REQ)
	condlog.Fatal(err, "Unable to create socket")
	defer sock.Close()

	err = sock.Connect(fmt.Sprintf("tcp://%s:%d", *hostPtr, *portPtr))
	condlog.Fatal(err, "Unable to connect")

	in := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for in.Scan() {
		hand := in.Text()

		results := solve(hand, sock)
		for _, r := range results {
			fmt.Println(r)
		}
		fmt.Print("> ")
	}

	condlog.Fatal(in.Err(), "Unable to read from stdin")
}
