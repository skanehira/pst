package gui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type ProcInfoView struct {
	*tview.TextView
}

func NewProcInfoView() *ProcInfoView {
	p := &ProcInfoView{
		TextView: tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true),
	}
	p.SetTitleAlign(tview.AlignLeft).SetTitle("process info").SetBorder(true)
	return p
}

func (p *ProcInfoView) UpdateInfo(g *Gui) {
	text := ""
	proc := g.ProcessManager.Selected()
	if proc != nil {
		info, err := g.ProcessManager.Info(proc.Pid)
		if err != nil {
			text = err.Error()
		} else {
			rows := strings.Split(info, "\n")
			if len(rows) == 1 {
				text = rows[0]
			} else {
				header := fmt.Sprintf("[yellow]%s[white]\n", rows[0])
				text = header + rows[1]
			}
		}
	}

	g.App.QueueUpdateDraw(func() {
		p.SetText(text)
	})
}

func (p *ProcInfoView) UpdateInfoWithPid(g *Gui, pid int) {
	text := ""
	if pid != 0 {
		info, err := g.ProcessManager.Info(pid)
		if err != nil {
			text = err.Error()
		} else {
			rows := strings.Split(info, "\n")
			if len(rows) == 1 {
				text = rows[0]
			} else {
				header := fmt.Sprintf("[yellow]%s[white]\n", rows[0])
				text = header + rows[1]
			}
		}
	}

	g.App.QueueUpdateDraw(func() {
		p.SetText(text)
	})

}
