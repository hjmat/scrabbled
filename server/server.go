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
        "github.com/hjmat/scrabbled/solver"
        "github.com/hjmat/scrabbled/logutil"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"flag"
	"fmt"
	"log"
        "os"
)

func serve(sock *zmq.Socket, solv *solver.Solver) {
     for {
          requestMsg, err := sock.Recv(0)
          logutil.Fatal("Unable to receive request", err)

          request := &scrabble.Request{}
          err = proto.Unmarshal([]byte(requestMsg), request)
          logutil.Fatal("Unable to unmarshal request", err)

          result := solv.Solve(*request.Hand)

          response := &scrabble.Response{Options: result}
          responseMsg, err := proto.Marshal(response)
          logutil.Fatal("Unable to marshal response", err)

          _, err = sock.Send(string(responseMsg), 0)
          logutil.Fatal("Unable to send response", err)
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

        solv := solver.NewSolver()

	err := solv.Populate(flag.Arg(0))
        logutil.Fatal("Unable to process corpus", err)

	sock, err := zmq.NewSocket(zmq.REP)
        logutil.Fatal("Unable to create socket", err)
        defer sock.Close()

	err = sock.Bind(fmt.Sprintf("tcp://*:%d", *portPtr))
	logutil.Fatal("Unable to bind socket", err)

	log.Printf("Listening on port %d", *portPtr)

        serve(sock, solv)
}
