package gui

import (
	"github.com/gdamore/tcell"
)

func (g *Gui) ProcessManagerKeybinds() {
	g.ProcessManager.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			g.SwitchPanel(g.FilterInput)
		case tcell.KeyF1:
			g.Help("process", g.ProcessManager)
		}

		switch event.Rune() {
		case 'K':
			if g.ProcessManager.Selected() != nil {
				g.Confirm("Do you want to kill this process?", "kill", g.ProcessManager, func() {
					g.ProcessManager.Kill()
					g.ProcessManager.UpdateView()
				})
			}
		}

		return event
	})

	g.ProcessManager.SetSelectionChangedFunc(func(row, col int) {
		if row < 1 {
			return
		}
		g.ProcInfoView.UpdateInfo(g)
		g.ProcessTreeView.UpdateTree(g)
	})
}

func (g *Gui) FilterInputKeybinds() {
	g.FilterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		case tcell.KeyEnter:
			g.SwitchPanel(g.ProcessManager)
			g.ProcInfoView.UpdateInfo(g)
			g.ProcessTreeView.UpdateTree(g)
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			g.SwitchPanel(g.ProcessManager)
			g.ProcInfoView.UpdateInfo(g)
			g.ProcessTreeView.UpdateTree(g)
		}
		return event
	})

	g.FilterInput.SetChangedFunc(func(text string) {
		g.ProcessManager.FilterWord = text
		g.ProcessManager.UpdateView()
	})
}

func (g *Gui) SetKeybinds() {
	g.FilterInputKeybinds()
	g.ProcessManagerKeybinds()
}
