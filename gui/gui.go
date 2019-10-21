package gui

import (
	"log"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Gui struct {
	FilterInput     *tview.InputField
	ProcessManager  *ProcessManager
	ProcessInfoView *ProcessInfoView
	ProcessTreeView *ProcessTreeView
	ProcessEnvView  *ProcessEnvView
	App             *tview.Application
	Pages           *tview.Pages
	Panels
}

type Panels struct {
	Current int
	Panels  []tview.Primitive
}

func New() *Gui {
	filterInput := tview.NewInputField().SetLabel("cmd name:")
	processManager := NewProcessManager()
	processInfoView := NewProcessInfoView()
	processTreeView := NewProcessTreeView(processManager)
	processEnvView := NewProcessEnvView()

	g := &Gui{
		FilterInput:     filterInput,
		ProcessManager:  processManager,
		App:             tview.NewApplication(),
		ProcessInfoView: processInfoView,
		ProcessTreeView: processTreeView,
		ProcessEnvView:  processEnvView,
	}

	g.Panels = Panels{
		Panels: []tview.Primitive{
			filterInput,
			processManager,
			processEnvView,
			processTreeView,
		},
	}

	return g
}

func (g *Gui) Confirm(message, doneLabel string, primitive tview.Primitive, doneFunc func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{doneLabel, "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.CloseAndSwitchPanel("modal", primitive)
			if buttonLabel == doneLabel {
				g.App.QueueUpdateDraw(func() {
					doneFunc()
				})
			}
		})

	g.Pages.AddAndSwitchToPage("modal", g.Modal(modal, 50, 29), true).ShowPage("main")
}

func (g *Gui) CloseAndSwitchPanel(removePrimitive string, primitive tview.Primitive) {
	g.Pages.RemovePage(removePrimitive).ShowPage("main")
	g.SwitchPanel(primitive)
}

func (g *Gui) TextView(title, text string, page tview.Primitive) {
	primitiveName := "textView"
	view := tview.NewTextView()
	view.SetTitle(title).SetTitleAlign(tview.AlignLeft).SetBorder(true)
	view.SetText(text)
	view.SetDynamicColors(true)

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			g.CloseAndSwitchPanel(primitiveName, page)
		}
		return event
	})

	g.Pages.AddAndSwitchToPage(primitiveName, g.Modal(view, 50, 20), true).ShowPage("main")
}

func (g *Gui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (g *Gui) SwitchPanel(p tview.Primitive) *tview.Application {
	g.ProcessInfoView.UpdateInfo(g)
	g.ProcessTreeView.UpdateTree(g)
	g.ProcessEnvView.UpdateView(g)
	return g.App.SetFocus(p)
}

func (g *Gui) Run() error {
	g.SetKeybinds()
	if err := g.ProcessManager.UpdateView(); err != nil {
		return err
	}
	// when start app, set select index 0
	g.ProcessManager.Select(1, 0)

	g.ProcessInfoView.UpdateInfo(g)
	g.ProcessTreeView.UpdateTree(g)
	g.ProcessEnvView.UpdateView(g)

	infoGrid := tview.NewGrid().SetRows(0, 0, 0).
		AddItem(g.ProcessInfoView, 0, 0, 1, 1, 0, 0, true).
		AddItem(g.ProcessEnvView, 1, 0, 1, 1, 0, 0, true).
		AddItem(g.ProcessTreeView, 2, 0, 1, 1, 0, 0, true)

	grid := tview.NewGrid().SetRows(1, 0).
		SetColumns(30, 0).
		AddItem(g.FilterInput, 0, 0, 1, 1, 0, 0, true).
		AddItem(g.ProcessManager, 1, 0, 1, 1, 0, 0, true).
		AddItem(infoGrid, 1, 1, 1, 1, 0, 0, true)

	g.Pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	if err := g.App.SetRoot(g.Pages, true).Run(); err != nil {
		g.App.Stop()
		log.Println(err)
		return err
	}

	return nil
}
