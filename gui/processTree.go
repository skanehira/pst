package gui

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type ProcessTreeView struct {
	*tview.TreeView
}

func NewProcessTreeView() *ProcessTreeView {
	p := &ProcessTreeView{
		TreeView: tview.NewTreeView(),
	}

	p.SetBorder(true).SetTitle("process tree").SetTitleAlign(tview.AlignLeft)
	return p
}

func (p *ProcessTreeView) UpdateTree(g *Gui) {
	proc := g.ProcessManager.Selected()
	if proc == nil {
		return
	}

	pid := strconv.Itoa(proc.Pid)

	root := tview.NewTreeNode(pid).
		SetColor(tcell.ColorYellow)

	p.SetRoot(root).
		SetCurrentNode(root)

	p.addNode(g, root, proc.Pid)

	p.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			pid := reference.(int)
			p.addNode(g, node, pid)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		}
		g.GrobalKeybind(event)
		return event
	})
}

func (p *ProcessTreeView) addNode(g *Gui, target *tview.TreeNode, pid int) {
	processes, err := g.ProcessManager.GetProcesses()
	if err != nil {
		return
	}

	proc, ok := processes[pid]
	if !ok {
		return
	}

	for _, p := range proc.Child {
		node := tview.NewTreeNode(fmt.Sprintf("PID: %d CMD: %s", p.Pid, p.Cmd)).
			SetReference(p.Pid)

		p, ok := processes[p.Pid]
		node.SetSelectable(ok)
		if len(p.Child) > 0 {
			node.SetColor(tcell.ColorGreen)
		}
		target.AddChild(node)
	}
}
