/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */
package keyloader

import (
	"errors"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"io/ioutil"
	"strings"
)

/*
 * Loads a zmq curve key pair from disk. If the requested pair cannot be found, a new pair that
 * matches the name is generated.
 */
func InitKeys(private_key_path string) (private_key string, public_key string, err error) {
	if !strings.HasSuffix(private_key_path, ".private") {
		return "", "", errors.New(fmt.Sprintf("Invalid private key path '%v'.", private_key_path))
	}

	public_key_path := strings.TrimSuffix(private_key_path, ".private") + ".public"

	private_key_buf, _ := ioutil.ReadFile(private_key_path)
	private_key = string(private_key_buf)
	public_key_buf, _ := ioutil.ReadFile(public_key_path)
	public_key = string(public_key_buf)

	if private_key != "" && public_key != "" {
		return private_key, public_key, nil
	} else if private_key == "" && public_key == "" {
		if public_key, private_key, err = zmq.NewCurveKeypair(); err != nil {
			return "", "", err
		}
		if err = ioutil.WriteFile(private_key_path, []byte(private_key), 0600); err != nil {
			return "", "", err
		}
		if err = ioutil.WriteFile(public_key_path, []byte(public_key), 0600); err != nil {
			return "", "", err
		}
	} else {
		return "", "", errors.New(fmt.Sprintf("Imbalanced key setup: '%v'", private_key_path))
	}

	return private_key, public_key, nil
}
