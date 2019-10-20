package gui

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Gui struct {
	FilterInput    *tview.InputField
	ProcessManager *ProcessManager
	App            *tview.Application
}

func New() *Gui {
	return &Gui{
		FilterInput:    tview.NewInputField(),
		ProcessManager: NewProcessManager(),
		App:            tview.NewApplication(),
	}
}

func (g *Gui) FilterInputKeybinds() {
	g.FilterInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		case tcell.KeyEnter:
			g.App.SetFocus(g.ProcessManager)
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			g.App.SetFocus(g.ProcessManager)
		}
		return event
	})

	g.FilterInput.SetChangedFunc(func(text string) {
		g.ProcessManager.FilterWord = text
		g.ProcessManager.UpdateView()
	})
}

func (g *Gui) ProcessManagerKeybinds() {
	g.ProcessManager.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEscape:
			g.App.Stop()
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			g.App.SetFocus(g.FilterInput)
		}

		return event
	})
}

func (g *Gui) SetKeybinds() {
	g.FilterInputKeybinds()
	g.ProcessManagerKeybinds()
}

func (g *Gui) Run() error {
	g.SetKeybinds()
	if err := g.ProcessManager.UpdateView(); err != nil {
		return err
	}

	grid := tview.NewGrid().SetRows(1, 0)
	grid.AddItem(g.FilterInput, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(g.ProcessManager, 1, 0, 2, 2, 0, 0, true)

	if err := g.App.SetRoot(grid, true).SetFocus(g.FilterInput).Run(); err != nil {
		g.App.Stop()
		log.Println(err)
		return err
	}

	return nil
}
