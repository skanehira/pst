package gui

import "github.com/rivo/tview"

type ProcInfoView struct {
	*tview.TextView
}

func NewProcInfoView() *ProcInfoView {
	p := &ProcInfoView{
		TextView: tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true),
	}
	p.SetBorder(true)
	return p
}

func (p *ProcInfoView) UpdateInfo(g *Gui, text string) {
	g.App.QueueUpdateDraw(func() {
		p.SetText(text)
	})
}
