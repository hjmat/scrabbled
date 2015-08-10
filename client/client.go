/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

/*
 * scrabbled client
 * USAGE: client [--port <PORT> --host <HOSTNAME>]
 */

package main

import (
	scrabble "github.com/hjmat/scrabbled/proto"
        "github.com/hjmat/scrabbled/logutil"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"bufio"
	"flag"
	"fmt"
	"os"
)

func solve(hand string, sock *zmq.Socket) []string {
	req := &scrabble.Request{Hand: proto.String(hand)}

	reqMsg, err := proto.Marshal(req)
        logutil.Fatal("Unable to marshal message", err)

	_, err = sock.Send(string(reqMsg), 0)
        logutil.Fatal("Unable to marshal message", err)

	resMsg, err := sock.Recv(0)
        logutil.Fatal("Unable to receive response", err)

	res := &scrabble.Response{}
	err = proto.Unmarshal([]byte(resMsg), res)
        logutil.Fatal("Unable to unmarshal request: ", err)

	return res.Options
}

func main() {
	hostPtr := flag.String("host", "localhost", "host")
	portPtr := flag.Int("port", 30000, "port")
	flag.Parse()

	sock, err := zmq.NewSocket(zmq.REQ)
        logutil.Fatal("Unable to create socket: ", err)
	defer sock.Close()

	err = sock.Connect(fmt.Sprintf("tcp://%s:%d", *hostPtr, *portPtr))
        logutil.Fatal("Unable to connect: ", err)

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

        logutil.Fatal("Unable to read from stdin: ", in.Err())
}
