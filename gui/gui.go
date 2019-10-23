package gui

import (
	"log"

	"github.com/rivo/tview"
)

const (
	InputPanel int = iota + 1
	ProcessesPanel
	ProcessInfoPanel
	ProcessEnvPanel
	ProcessTreePanel
)

type Gui struct {
	FilterInput     *tview.InputField
	ProcessManager  *ProcessManager
	ProcessInfoView *ProcessInfoView
	ProcessTreeView *ProcessTreeView
	ProcessEnvView  *ProcessEnvView
	NaviView        *NaviView
	App             *tview.Application
	Pages           *tview.Pages
	Panels
}

type Panels struct {
	Current int
	Panels  []tview.Primitive
	Kinds   []int
}

func New() *Gui {
	filterInput := tview.NewInputField().SetLabel("cmd name:")
	processManager := NewProcessManager()
	processInfoView := NewProcessInfoView()
	processTreeView := NewProcessTreeView(processManager)
	processEnvView := NewProcessEnvView()
	naviView := NewNaviView()

	g := &Gui{
		FilterInput:     filterInput,
		ProcessManager:  processManager,
		App:             tview.NewApplication(),
		ProcessInfoView: processInfoView,
		ProcessTreeView: processTreeView,
		ProcessEnvView:  processEnvView,
		NaviView:        naviView,
	}

	g.Panels = Panels{
		Panels: []tview.Primitive{
			filterInput,
			processManager,
			processInfoView,
			processEnvView,
			processTreeView,
		},
		Kinds: []int{
			InputPanel,
			ProcessesPanel,
			ProcessInfoPanel,
			ProcessEnvPanel,
			ProcessTreePanel,
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
	g.NaviView.UpdateView(g)
	return g.App.SetFocus(p)
}

func (g *Gui) CurrentPanelKind() int {
	return g.Panels.Kinds[g.Panels.Current]
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
	g.NaviView.UpdateView(g)

	infoGrid := tview.NewGrid().SetRows(0, 0, 0).
		SetColumns(30, 0).
		AddItem(g.ProcessManager, 0, 0, 3, 1, 0, 0, true).
		AddItem(g.ProcessInfoView, 0, 1, 1, 1, 0, 0, true).
		AddItem(g.ProcessEnvView, 1, 1, 1, 1, 0, 0, true).
		AddItem(g.ProcessTreeView, 2, 1, 1, 1, 0, 0, true)

	grid := tview.NewGrid().SetRows(1, 0, 2).
		SetColumns(30).
		AddItem(g.FilterInput, 0, 0, 1, 1, 0, 0, true).
		AddItem(infoGrid, 1, 0, 1, 2, 0, 0, true).
		AddItem(g.NaviView, 2, 0, 1, 2, 0, 0, true)

	g.Pages = tview.NewPages().
		AddAndSwitchToPage("main", grid, true)

	if err := g.App.SetRoot(g.Pages, true).Run(); err != nil {
		g.App.Stop()
		log.Println(err)
		return err
	}

	return nil
}
