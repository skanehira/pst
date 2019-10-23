package gui

import (
	"github.com/rivo/tview"
)

type ProcessFileView struct {
	*tview.TextView
}

func NewProcessFileView() *ProcessFileView {
	p := &ProcessFileView{
		TextView: tview.NewTextView().SetDynamicColors(true),
	}

	p.SetTitleAlign(tview.AlignLeft).SetTitle("process open files").SetBorder(true)
	p.SetWrap(false)
	return p
}

func (p *ProcessFileView) UpdateViewWithPid(g *Gui, pid int) {
	text := ""
	if pid != 0 {
		info, err := g.ProcessManager.OpenFiles(pid)
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

func (p *ProcessFileView) UpdateView(g *Gui) {
	proc := g.ProcessManager.Selected()
	if proc != nil {
		p.UpdateViewWithPid(g, proc.Pid)
	}
}
