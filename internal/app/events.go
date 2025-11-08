package app

import (
	"fmt"

	"github.com/rivo/tview"
)

func (a *App) generateEventsList(listEvents *tview.List) {
	listEvents.Clear()

	events := a.eventsSnapshot()

	if len(events) == 0 {
		listEvents.AddItem("No events yet!", "", 0, nil)
		listEvents.AddItem("Start interacting with your tamagotchi to see events.", "", 0, nil)
		return
	}

	listEvents.AddItem("=== GAME EVENTS ===", "", 0, nil)
	listEvents.AddItem("", "", 0, nil) // Empty line

	// Show events in reverse chronological order (newest first)
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		timeStr := event.Timestamp.Format("15:04")

		// Add color coding based on event type
		var eventIcon string
		switch event.Type {
		case "FEED":
			eventIcon = "🍽️"
		case "PLAY":
			eventIcon = "🎮"
		case "SLEEP":
			eventIcon = "😴"
		case "EVOLUTION":
			eventIcon = "🌟"
		case "DEATH":
			eventIcon = "💔"
		case "RESTART":
			eventIcon = "🔄"
		default:
			eventIcon = "📝"
		}

		listEvents.AddItem(
			fmt.Sprintf("%s [%s] %s", eventIcon, timeStr, event.Message),
			"",
			0,
			nil,
		)
	}
}

func (a *App) eventsPage() (title string, content tview.Primitive) {
	listEvents := a.viewsList["events"]
	if listEvents == nil {
		listEvents = getList()
		a.viewsList["events"] = listEvents
	}

	a.generateEventsList(listEvents)

	title = eventsSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listEvents, 0, 1, true), 0, 1, true)
}
