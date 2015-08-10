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
	"github.com/hjmat/scrabbled/condlog"
	scrabble "github.com/hjmat/scrabbled/proto"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"flag"
	"fmt"
	"log"
	"os"
)

func serve(sock *zmq.Socket, solv *Solver) {
	for {
		requestMsg, err := sock.Recv(0)
		condlog.Fatal(err, "Unable to receive request")

		request := &scrabble.Request{}
		err = proto.Unmarshal([]byte(requestMsg), request)
		condlog.Fatal(err, "Unable to unmarshal request")

		result := solv.Solve(request.Hand)

		response := &scrabble.Response{Words: result}
		responseMsg, err := proto.Marshal(response)
		condlog.Fatal(err, "Unable to marshal response")

		_, err = sock.Send(string(responseMsg), 0)
		condlog.Fatal(err, "Unable to send response")
	}
}

func usage() {
	fmt.Println("Usage: server [OPTIONS] <PATH TO CORPUS>")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	portPtr := flag.Int("port", 30000, "port")
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}

	solv := NewSolver()

	err := solv.Populate(flag.Arg(0))
	condlog.Fatal(err, "Unable to process corpus")

	sock, err := zmq.NewSocket(zmq.REP)
	condlog.Fatal(err, "Unable to create socket")
	defer sock.Close()

	err = sock.Bind(fmt.Sprintf("tcp://*:%d", *portPtr))
	condlog.Fatal(err, "Unable to bind socket")

	log.Printf("Listening on port %d", *portPtr)

	serve(sock, solv)
}
