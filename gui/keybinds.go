package gui

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (g *Gui) nextPanel() {
	idx := (g.Panels.Current + 1) % len(g.Panels.Panels)
	g.Panels.Current = idx
	g.SwitchPanel(g.Panels.Panels[g.Panels.Current])
}

func (g *Gui) prePanel() {
	g.Panels.Current--

	if g.Panels.Current < 0 {
		g.Current = len(g.Panels.Panels) - 1
	} else {
		idx := (g.Panels.Current) % len(g.Panels.Panels)
		g.Panels.Current = idx
	}
	g.SwitchPanel(g.Panels.Panels[g.Panels.Current])
}

func (g *Gui) GrobalKeybind(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyTab:
		g.nextPanel()
	case tcell.KeyBacktab:
		g.prePanel()
	}

	g.NaviView.UpdateView(g)
}

func (g *Gui) ProcessManagerKeybinds() {
	g.ProcessManager.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'K':
			if g.ProcessManager.Selected() != nil {
				g.Confirm("Do you want to kill this process?", "kill", g.ProcessManager, func() {
					g.ProcessManager.Kill()
					g.ProcessManager.UpdateView()
				})
			}
		}

		g.GrobalKeybind(event)
		return event
	})

	g.ProcessManager.SetSelectionChangedFunc(func(row, col int) {
		if row < 1 {
			return
		}
		g.ProcessInfoView.UpdateInfo(g)
		g.ProcessTreeView.UpdateTree(g)
		g.ProcessEnvView.UpdateView(g)
		g.ProcessFileView.UpdateView(g)
	})
}

func (g *Gui) FilterInputKeybinds() {
	g.FilterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		case tcell.KeyEnter:
			g.nextPanel()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.GrobalKeybind(event)
		return event
	})

	g.FilterInput.SetChangedFunc(func(text string) {
		g.ProcessManager.FilterWord = text
		g.ProcessManager.UpdateView()
	})
}

func (g *Gui) ProcessTreeViewKeybinds() {
	g.ProcessTreeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'K':
			if node := g.ProcessTreeView.GetCurrentNode(); node != nil {
				if ref := node.GetReference(); ref != nil {
					g.Confirm("Do you want to kill this process?", "kill", g.ProcessTreeView, func() {
						g.ProcessManager.KillWithPid(ref.(int))
						// wait a little to finish process killing
						time.Sleep(1 * time.Millisecond)
						g.ProcessTreeView.UpdateTree(g)
					})
				}
			}
		}
		g.GrobalKeybind(event)
		return event
	})

	g.ProcessTreeView.SetChangedFunc(func(node *tview.TreeNode) {
		if node == nil {
			return
		}
		ref := node.GetReference()
		if ref == nil {
			return
		}

		pid := ref.(int)
		g.ProcessInfoView.UpdateInfoWithPid(g, pid)
		g.ProcessEnvView.UpdateViewWithPid(g, pid)
		g.ProcessFileView.UpdateViewWithPid(g, pid)
	})
}

func (g *Gui) ProcessEnvViewKeybinds() {
	g.ProcessEnvView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.GrobalKeybind(event)
		return event
	})
}

func (g *Gui) ProcessInfoViewKeybinds() {
	g.ProcessInfoView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.GrobalKeybind(event)
		return event
	})
}

func (g *Gui) ProcessFileViewKeybinds() {
	g.ProcessFileView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.GrobalKeybind(event)
		return event
	})
}

func (g *Gui) SetKeybinds() {
	g.FilterInputKeybinds()
	g.ProcessManagerKeybinds()
	g.ProcessTreeViewKeybinds()
	g.ProcessInfoViewKeybinds()
	g.ProcessEnvViewKeybinds()
	g.ProcessFileViewKeybinds()
}
