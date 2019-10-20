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

func NewProcessTreeView(pm *ProcessManager) *ProcessTreeView {
	p := &ProcessTreeView{
		TreeView: tview.NewTreeView(),
	}

	p.SetBorder(true).SetTitle("process tree").SetTitleAlign(tview.AlignLeft)

	p.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			pid := reference.(int)
			p.addNode(pm, node, pid)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
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

	p.addNode(g.ProcessManager, root, proc.Pid)
}

func (p *ProcessTreeView) addNode(pm *ProcessManager, target *tview.TreeNode, pid int) {
	processes, err := pm.GetProcesses()
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
