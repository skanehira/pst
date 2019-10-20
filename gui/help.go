package gui

import "github.com/rivo/tview"

var helps = map[string]string{
	"process": `[yellow]Keybinds[white]
[red]j[white]	next process
[red]k[white]	previous process
[red]g[white]	first process
[red]G[white]	last process
[red]K[white]	kill selected process
`,
	"filter": `
`,
}

func (g *Gui) Help(key string, page tview.Primitive) {
	g.TextView("HELP", helps[key], page)
}
