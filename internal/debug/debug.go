package debug

import (
	"log"
)

func Logf(debug bool, message string, arg string) {
	if debug {
		log.Printf(message, arg)
	}
}

func Logln(debug bool, message string) {
	if debug {
		log.Println(message)
	}
}
