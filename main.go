package main

import (
	"os"

	"github.com/regalias/atlas-app/appserver"
)

func main() {
	os.Exit(appserver.Run(os.Args[1:]))
}
