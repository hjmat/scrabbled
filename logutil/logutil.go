/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

package logutil

import (
       	"log"
)

func Error(msg string, err error) {
     	if err != nil {
		log.Print(msg, ": ", err)
	}
}

func Fatal(msg string, err error) {
	if err != nil {
		log.Fatal(msg, ": ", err)
	}
}
