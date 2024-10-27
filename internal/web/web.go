package web

import (
	"github.com/calebbramel/azpgp/internal/azenv"
)

type RequestBody struct {
	Recipient string `json:"username"`
	ID        string `json:"id"`
}

type Response struct {
	Message string `json:"message"`
}

func init() {
	azenv.Load()
}
