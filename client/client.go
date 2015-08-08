/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

package main

import (
	scrabble "github.com/hjmat/scrabbled/proto"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func solve(hand string, sock *zmq.Socket) []string {
	req := &scrabble.Request{Hand: proto.String(hand)}

	reqMsg, err := proto.Marshal(req)
	if err != nil {
		log.Fatal("Unable to marshal message: ", err)
	}

	_, err = sock.Send(string(reqMsg), 0)
	if err != nil {
		log.Fatal("Unable to send request: ", err)
	}

	resMsg, err := sock.Recv(0)
	if err != nil {
		log.Fatal("Unable to receive response: ", err)
	}

	res := &scrabble.Response{}
	err = proto.Unmarshal([]byte(resMsg), res)
	if err != nil {
		log.Fatal("Unable to unmarshal request: ", err)
	}

	return res.Options
}

func main() {
	hostPtr := flag.String("host", "localhost", "host")
	portPtr := flag.Int("port", 30000, "port")
	flag.Parse()

	sock, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		log.Fatal("Unable to create socket: ", err)
	}
	defer sock.Close()

	err = sock.Connect(fmt.Sprintf("tcp://%s:%d", *hostPtr, *portPtr))
	if err != nil {
		log.Fatal("Unable to connect: ", err)
	}

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

	if in.Err() != nil {
		log.Fatal("Unable to read from stdin: ", in.Err())
	}
}
