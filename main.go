package main

import (
	"fmt"
	"log"

	ps "github.com/mitchellh/go-ps"
)

func main() {
	processes, err := ps.Processes()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range processes {
		fmt.Printf("%#+v\n", p)
	}
}
