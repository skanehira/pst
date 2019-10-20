package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/skanehira/pst/gui"
)

var (
	enableLog = flag.Bool("log", false, "enable output log")
)

func run() int {
	flag.Parse()

	if *enableLog {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		logWriter, err := os.OpenFile(filepath.Join(home, "pst.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		log.SetOutput(logWriter)
		log.SetFlags(log.Lshortfile)
	}

	if err := gui.New().Run(); err != nil {
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
