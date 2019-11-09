package gui

import (
	"fmt"

	"github.com/rivo/tview"
)

type NaviView struct {
	*tview.TextView
}

func NewNaviView() *NaviView {
	n := &NaviView{
		TextView: tview.NewTextView().SetTextAlign(tview.AlignLeft).SetDynamicColors(true),
	}
	n.SetTitleAlign(tview.AlignLeft)
	return n
}

func (n *NaviView) UpdateView(g *Gui) {
	g.App.QueueUpdateDraw(func() {
		switch g.CurrentPanelKind() {
		case InputPanel:
			n.SetText(fmt.Sprintf("%s", switchNavi))
		case ProcessesPanel:
			n.SetText(fmt.Sprintf("%s, %s, %s", moveNavi, switchNavi, helps[ProcessesPanel]))
		case ProcessInfoPanel:
			n.SetText(fmt.Sprintf("%s, %s, %s", moveNavi, switchNavi, helps[ProcessInfoPanel]))
		case ProcessEnvPanel:
			n.SetText(fmt.Sprintf("%s, %s, %s", moveNavi, switchNavi, helps[ProcessEnvPanel]))
		case ProcessTreePanel:
			n.SetText(fmt.Sprintf("%s, %s, %s", moveNavi, switchNavi, helps[ProcessTreePanel]))
		case ProcessFilePanel:
			n.SetText(fmt.Sprintf("%s, %s, %s", moveNavi, switchNavi, helps[ProcessFilePanel]))
		default:
			n.SetText("")
		}
	})
}

var (
	moveNavi   = "[red::b]j[white]: move down, [red]k[white]: move up, [red]h[white]: move left, [red]l[white]: move right, [red]g[white]: move to top, [red]G[white]: move to bottom, [red]Ctrl-f[white]: next page [red]Ctrl-b[white]: previous page, [red]Ctrl-c[white]: stop pst"
	switchNavi = `[red::b]Tab[white]: next panel, [red]Shift-Tab[white]: previous panel`
)

var helps = map[int]string{
	InputPanel:       ``,
	ProcessesPanel:   `[red]K[white]: kill process`,
	ProcessInfoPanel: ``,
	ProcessEnvPanel:  ``,
	ProcessTreePanel: `[red]K[white]: kill process, [red]h[white]: collapse, [red]l[white]: expand, [red]enter[white]: expand toggle`,
	ProcessFilePanel: ``,
}
