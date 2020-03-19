package main

import (
	"os"

	"github.com/regalias/atlas-app/server"
)

func main() {
	os.Exit(server.Run(os.Args[1:]))
}
