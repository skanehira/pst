package gui

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	ps "github.com/mitchellh/go-ps"
	"github.com/rivo/tview"
)

type ProcessManager struct {
	*tview.Table
	processes  []Process
	FilterWord string
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		Table: tview.NewTable().Select(0, 0).SetFixed(1, 1).SetSelectable(true, false),
	}
}

func (p *ProcessManager) GetProcesses() error {
	processes, err := ps.Processes()
	if err != nil {
		log.Println("cannot get processes: " + err.Error())
		return err
	}

	pids := make(map[int]Process)
	for _, proc := range processes {
		// skip pid 0
		if proc.Pid() == 0 {
			continue
		}
		if strings.Index(proc.Executable(), p.FilterWord) == -1 {
			continue
		}
		pids[proc.Pid()] = Process{
			Pid:  proc.Pid(),
			PPid: proc.PPid(),
			Cmd:  proc.Executable(),
		}
	}

	// add child processes
	for _, p := range processes {
		if p.Pid() == p.PPid() {
			continue
		}

		if proc, ok := pids[p.PPid()]; ok {
			proc.Child = append(proc.Child, pids[p.Pid()])
			pids[p.PPid()] = proc
		}
	}

	p.processes = []Process{}
	for _, proc := range pids {
		p.processes = append(p.processes, proc)
	}

	sort.Slice(p.processes, func(i, j int) bool {
		return p.processes[i].Pid < p.processes[j].Pid
	})

	return nil
}

var headers = []string{
	"Pid",
	"PPid",
	"Cmd",
}

func (p *ProcessManager) UpdateView() error {
	// get processes
	if err := p.GetProcesses(); err != nil {
		return err
	}

	table := p.Clear()

	// set headers
	for i, h := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            h,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorYellow,
			BackgroundColor: tcell.ColorDefault,
		})
	}

	// set process info to cell
	var i int
	for _, proc := range p.processes {
		pid := strconv.Itoa(proc.Pid)
		ppid := strconv.Itoa(proc.PPid)
		table.SetCell(i+1, 0, tview.NewTableCell(pid))
		table.SetCell(i+1, 1, tview.NewTableCell(ppid))
		table.SetCell(i+1, 2, tview.NewTableCell(proc.Cmd))
		i++
	}

	return nil
}

func (p *ProcessManager) Selected() *Process {
	if len(p.processes) == 0 {
		return nil
	}
	row, _ := p.GetSelection()
	if row < 0 {
		return nil
	}
	return &p.processes[row-1]
}

func (p *ProcessManager) Kill() error {
	pid := p.Selected().Pid
	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := proc.Kill(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *ProcessManager) Info() (string, error) {
	// TODO implements windows
	if runtime.GOOS == "windows" {
		return "", nil
	}

	proc := p.Selected()
	if proc == nil {
		return "", nil
	}
	buf := bytes.Buffer{}
	cmd := exec.Command("ps", "-o", "pid,ppid,%cpu,%mem,lstart,user,command", "-p", strconv.Itoa(proc.Pid))
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
