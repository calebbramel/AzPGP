package logger

import (
	"log"
)

func HandleErrf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func HandleErrln(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
