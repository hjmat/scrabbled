/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

/*
 * scrabbled server
 */

package main

import (
	"github.com/hjmat/scrabbled/condlog"
	"github.com/hjmat/scrabbled/keyloader"
	scrabble "github.com/hjmat/scrabbled/proto"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
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

func initSecurity(client_key_path string, private_key_path string, sock *zmq.Socket) {
	zmq.AuthStart()
	private_key, _, err := keyloader.InitKeys(private_key_path)
	condlog.Fatal(err, fmt.Sprintf("Unable to read key pair for private key '%v'", private_key_path))
	sock.ServerAuthCurve("scrabble", private_key)

	// Add all the public keys in the client key directory
	files, err := ioutil.ReadDir(client_key_path)
	condlog.Fatal(err, fmt.Sprintf("Unable to enumerate client keys in '%v'", client_key_path))
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".public") {
			fullpath := path.Join(client_key_path, f.Name())
			err = keyloader.CheckPermissions(fullpath)
			condlog.Fatal(err, "Untrustworthy key file")
			buf, err := ioutil.ReadFile(fullpath)
			condlog.Fatal(err, fmt.Sprintf("Unable to load public client key '%v'", fullpath))
			zmq.AuthCurveAdd("scrabble", string(buf))
		}
	}
}

func main() {
	portPtr := flag.Int("port", 30000, "port")
	keyPtr := flag.String("key", "server.private", "path to private key")
	clientPtr := flag.String("clientkeys", "clients.allow", "path to directory containing public keys of clients")
	corpusPtr := flag.String("corpus", "corpus.txt", "path to a newline-separated list of words")
	flag.Parse()

	solv := NewSolver()
	err := solv.Populate(path.Clean(*corpusPtr))
	condlog.Fatal(err, "Unable to process corpus")

	sock, err := zmq.NewSocket(zmq.REP)
	condlog.Fatal(err, "Unable to create socket")
	defer sock.Close()
	initSecurity(path.Clean(*clientPtr), path.Clean(*keyPtr), sock)

	err = sock.Bind(fmt.Sprintf("tcp://*:%d", *portPtr))
	condlog.Fatal(err, "Unable to bind socket")

	log.Printf("Listening on port %d", *portPtr)

	serve(sock, solv)
}
