package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	statusSection = "Status"
	feedSection   = "Feed"
	playSection   = "Play"
	sleepSection  = "Sleep"
	eventsSection = "Events"
	helpSection   = "Help"
)

func (a *App) getPagesInfo() (tview.Primitive, tview.Primitive) {
	pages := tview.NewPages()

	// Add all pages
	_, statusContent := a.statusPage()
	pages.AddPage(statusSection, statusContent, true, true)

	_, feedContent := a.feedPage()
	pages.AddPage(feedSection, feedContent, true, false)

	_, playContent := a.playPage()
	pages.AddPage(playSection, playContent, true, false)

	_, sleepContent := a.sleepPage()
	pages.AddPage(sleepSection, sleepContent, true, false)

	_, eventsContent := a.eventsPage()
	pages.AddPage(eventsSection, eventsContent, true, false)

	_, helpContent := a.helpPage()
	pages.AddPage(helpSection, helpContent, true, false)

	// Info bar
	info := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Ctrl+S: Status | Ctrl+F: Feed | Ctrl+P: Play | Ctrl+L: Sleep | Ctrl+E: Events | Ctrl+H: Help | Ctrl+R: Restart | Ctrl+C: Quit")

	return pages, info
}

func (a *App) goToSection(section string, info tview.Primitive) {
	pages := a.TLayout.GetItem(0).(*tview.Pages)
	pages.SwitchToPage(section)
	pages.SendToFront(section)
}

func getList() *tview.List {
	list := tview.NewList()
	list.ShowSecondaryText(false)
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorYellow)
	list.SetBorder(true)
	list.SetTitle("Termagotchi")
	return list
}
