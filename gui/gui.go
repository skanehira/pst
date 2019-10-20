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
	Pages          *tview.Pages
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

		switch event.Rune() {
		case 'K':
			g.Confirm("Do you want to kill this process?", "kill", g.ProcessManager, func() {
				g.ProcessManager.Kill()
				g.ProcessManager.UpdateView()
			})
		}

		return event
	})
}

func (g *Gui) SetKeybinds() {
	g.FilterInputKeybinds()
	g.ProcessManagerKeybinds()
}

func (g *Gui) Confirm(message, doneLabel string, panel tview.Primitive, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{doneLabel, "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.CloseAndSwitchPanel("modal", panel)
			if buttonLabel == doneLabel {
				doneFunc()
			}
		})

	g.Pages.AddAndSwitchToPage("modal", g.Modal(modal, 50, 29), true).ShowPage("main")
}

func (g *Gui) CloseAndSwitchPanel(removePanel string, panel tview.Primitive) {
	g.Pages.RemovePage(removePanel).ShowPage("main")
	g.SwitchPanel(panel)
}

func (g *Gui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (g *Gui) SwitchPanel(p tview.Primitive) *tview.Application {
	return g.App.SetFocus(p)
}

func (g *Gui) Run() error {
	g.SetKeybinds()
	if err := g.ProcessManager.UpdateView(); err != nil {
		return err
	}

	grid := tview.NewGrid().SetRows(1).
		AddItem(g.FilterInput, 0, 0, 1, 1, 0, 0, true).
		AddItem(g.ProcessManager, 1, 0, 2, 1, 0, 0, true)

	g.Pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	if err := g.App.SetRoot(g.Pages, true).Run(); err != nil {
		g.App.Stop()
		log.Println(err)
		return err
	}

	return nil
}
