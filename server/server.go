/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

/*
 * scrabbled server
 * USAGE: server [--port <PORT>] <path to corpus file>
 */

package main

import (
	scrabble "github.com/hjmat/scrabbled/proto"
	solver "github.com/hjmat/scrabbled/solver"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"flag"
	"fmt"
	"log"
)

func fatal(msg string, err error) {
	if err != nil {
		log.Fatal(msg, ": ", err)
	}
}

func main() {
	portPtr := flag.Int("port", 30000, "port")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Usage: server [--port <port>] CORPUS")
	}

	err := solver.Populate(flag.Arg(0))
        fatal("Unable to process corpus", err)

	sock, err := zmq.NewSocket(zmq.REP)
	fatal("Unable to create socket", err)

	defer sock.Close()

	err = sock.Bind(fmt.Sprintf("tcp://*:%d", *portPtr))
	fatal("Unable to bind socket", err)

	log.Printf("Listening on port %d", *portPtr)

	for {
		requestMsg, err := sock.Recv(0)
		fatal("Unable to receive request", err)

		request := &scrabble.Request{}
		err = proto.Unmarshal([]byte(requestMsg), request)
		fatal("Unable to unmarshal request", err)

		result := solver.Solve(*request.Hand)

		reply := &scrabble.Response{Options: result}
		replyMsg, err := proto.Marshal(reply)
		fatal("Unable to marshal response", err)

		_, err = sock.Send(string(replyMsg), 0)
		fatal("Unable to send response", err)
	}
}
