/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

/*
 * scrabbled client
 */

package main

import (
	"github.com/hjmat/scrabbled/condlog"
	"github.com/hjmat/scrabbled/keyloader"
	scrabble "github.com/hjmat/scrabbled/proto"

	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"

	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func initSecurity(private_key_path string, server_key_path string, sock *zmq.Socket) {
	zmq.AuthStart()
	private_key, public_key, err := keyloader.InitKeys(private_key_path)
	condlog.Fatal(err, fmt.Sprintf("Unable to read key pair for private key '%v'", private_key_path))
	zmq.AuthCurveAdd("scrabble", public_key)

	server_key_buf, err := ioutil.ReadFile(server_key_path)
	condlog.Fatal(err, fmt.Sprintf("Unable to load public server key '%v'", server_key_path))
	server_key := string(server_key_buf)
	sock.ClientAuthCurve(server_key, public_key, private_key)
}

func main() {
	hostPtr := flag.String("host", "localhost", "hostname of the server")
	portPtr := flag.Int("port", 30000, "port that the server runs on")
	keyPtr := flag.String("key", "client.private", "private key for authentication")
	servKeyPtr := flag.String("servkey", "server.public", "public key of the server")
	flag.Parse()

	sock, err := zmq.NewSocket(zmq.REQ)
	condlog.Fatal(err, "Unable to create socket")
	defer sock.Close()
	initSecurity(path.Clean(*keyPtr), path.Clean(*servKeyPtr), sock)

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
