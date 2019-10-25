package gui

import (
	"github.com/rivo/tview"
)

type ProcessEnvView struct {
	*tview.TextView
}

func NewProcessEnvView() *ProcessEnvView {
	p := &ProcessEnvView{
		TextView: tview.NewTextView().SetDynamicColors(true),
	}

	p.SetTitleAlign(tview.AlignLeft).SetTitle("process environments").SetBorder(true)
	p.SetWrap(false)
	return p
}

func (p *ProcessEnvView) UpdateViewWithPid(g *Gui, pid int) {
	text := ""
	if pid != 0 {
		info, err := g.ProcessManager.Env(pid)
		if err != nil {
			text = err.Error()
		} else {
			text = info
		}
	}

	g.App.QueueUpdateDraw(func() {
		p.SetText(text)
		p.ScrollToBeginning()
	})
}

func (p *ProcessEnvView) UpdateView(g *Gui) {
	proc := g.ProcessManager.Selected()
	if proc != nil {
		g.ProcessEnvView.UpdateViewWithPid(g, proc.Pid)
	}
}
