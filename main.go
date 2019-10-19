package main

import (
	"fmt"
	"log"
	"os"

	ps "github.com/mitchellh/go-ps"
)

func parseProcesses(processes []ps.Process) map[int]Process {
	// get ppid list
	pids := make(map[int]Process)
	for _, p := range processes {
		pids[p.Pid()] = Process{
			Pid:  p.Pid(),
			PPid: p.PPid(),
			Cmd:  p.Executable(),
		}
	}

	for _, p := range processes {
		if p.Pid() == p.PPid() {
			continue
		}

		if proc, ok := pids[p.PPid()]; ok {
			proc.Child = append(proc.Child, pids[p.Pid()])
			pids[p.PPid()] = proc
		}
	}

	return pids
}

func run() int {
	// get processes
	processes, err := ps.Processes()
	if err != nil {
		log.Println("cannot get processes: " + err.Error())
		return 1
	}

	for _, p := range parseProcesses(processes) {
		fmt.Printf("%#+v\n", p)
	}

	return 0
}

func main() {
	os.Exit(run())
}
