/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

package condlog

import (
	"log"
)

func Print(err error, msg ...string) {
	if err != nil {
		log.Print(msg, ": ", err)
	}
}

func Fatal(err error, msg ...string) {
	if err != nil {
		log.Fatal(msg, ": ", err)
	}
}
