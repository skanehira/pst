package gui

import (
	"bytes"
	"errors"
	"fmt"
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

var psArgs = GetEnv("PS_ARGS", "pid,ppid,%cpu,%mem,lstart,user,command")

type ProcessManager struct {
	*tview.Table
	processes  []Process
	FilterWord string
}

func NewProcessManager() *ProcessManager {
	p := &ProcessManager{
		Table: tview.NewTable().Select(0, 0).SetFixed(1, 1).SetSelectable(true, false),
	}
	p.SetBorder(true).SetTitle("processes").SetTitleAlign(tview.AlignLeft)
	return p
}

func (p *ProcessManager) GetProcesses() (map[int]Process, error) {
	processes, err := ps.Processes()
	if err != nil {
		log.Println("cannot get processes: " + err.Error())
		return nil, err
	}

	pids := make(map[int]Process)
	for _, proc := range processes {
		// skip pid 0
		if proc.Pid() == 0 {
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
		if strings.Index(proc.Cmd, p.FilterWord) == -1 {
			continue
		}

		p.processes = append(p.processes, proc)
	}

	sort.Slice(p.processes, func(i, j int) bool {
		return p.processes[i].Pid < p.processes[j].Pid
	})

	return pids, nil
}

var headers = []string{
	"Pid",
	"PPid",
	"Cmd",
}

func (p *ProcessManager) UpdateView() error {
	// get processes
	if _, err := p.GetProcesses(); err != nil {
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
	if len(p.processes) < row {
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

func (p *ProcessManager) KillWithPid(pid int) error {
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

func (p *ProcessManager) Info(pid int) (string, error) {
	// TODO implements windows
	if runtime.GOOS == "windows" {
		return "", nil
	}

	if pid == 0 {
		return "", nil
	}

	buf := bytes.Buffer{}
	cmd := exec.Command("ps", "-o", psArgs, "-p", strconv.Itoa(pid))
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", errors.New(buf.String())
	}

	return buf.String(), nil
}

func (p *ProcessManager) Env(pid int) (string, error) {
	// TODO implements windows
	if runtime.GOOS == "windows" {
		return "", nil
	}

	if pid == 0 {
		return "", nil
	}

	buf := bytes.Buffer{}
	cmd := exec.Command("ps", "eww", "-o", "command", "-p", strconv.Itoa(pid))
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.Split(buf.String(), "\n")

	var (
		envStr []string
		envs   []string
	)

	if len(result) > 1 {
		envStr = strings.Split(result[1], " ")[1:]
	} else {
		return buf.String(), nil
	}

	for _, e := range envStr {
		kv := strings.Split(e, "=")
		if len(kv) != 2 {
			continue
		}
		envs = append(envs, fmt.Sprintf("[yellow]%s[white]\t%s", kv[0], kv[1]))
	}

	return strings.Join(envs, "\n"), nil
}
