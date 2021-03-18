package main

import (
	"os"
	"strings"

	"github.com/ProfessorSnep/pubpipe/server"
)

func main() {
	port, exists := os.LookupEnv("PUBPIPE_PORT")
	if !exists {
		port = "8041"
	} else {
		port = strings.TrimSpace(port)
	}

	server.Start(port)
}
