package main

type Process struct {
	Pid   int
	PPid  int
	Cmd   string
	Child []Process
}
