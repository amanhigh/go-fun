package tui

import "github.com/gdamore/tcell/v2"

func (d *Darius) setupKeyBindings() {
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			d.handleEnterKey()
		case tcell.KeyEscape:
			d.handleEscapeKey()
		case tcell.KeyRune:
			d.handleRuneKey(event.Rune())
		}
		return event
	})
}

func (d *Darius) handleEnterKey() {
	if d.app.GetFocus() == d.availableServices {
		d.toggleServiceSelection(d.availableServices.GetCurrentItem(), "", "", 0)
	} else if d.app.GetFocus() == d.filterInput {
		d.applyFilterToSelection()
	}
}

func (d *Darius) handleEscapeKey() {
	if d.app.GetFocus() == d.filterInput {
		d.filterInput.SetText("")
		d.filterServices("")
		d.app.SetFocus(d.availableServices)
	}
}

func (d *Darius) handleRuneKey(r rune) {
	switch r {
	case 'i', 'I':
		d.executeCommand(CommandSetup)
	case 'c', 'C':
		d.executeCommand(CommandClean)
	case 'm':
		d.showMinikubePanel()
	case '/':
		d.app.SetFocus(d.filterInput)
	case 'q', 'Q':
		d.app.Stop()
	}
}
