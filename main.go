package main

import (
	"log"

	"github.com/ootiny/gapi/core"
)

func main() {
	if err := core.Output(); err != nil {
		log.Panicf("Failed to output: %v", err)
	}
}
