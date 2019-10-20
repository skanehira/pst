package gui

type Process struct {
	Pid   int
	PPid  int
	Cmd   string
	Child []Process
}
