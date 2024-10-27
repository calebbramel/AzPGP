package logger

import (
	"log"

	"github.com/calebbramel/azpgp/internal/azenv"
)

func Debugf(message string, arg interface{}) {
	if azenv.DebugFlag {
		log.Printf(message, arg)
	}
}

func Debugln(message string) {
	if azenv.DebugFlag {
		log.Println(message)
	}
}
