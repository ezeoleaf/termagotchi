package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func (a *App) generateStatusList(listStatus *tview.List) {
	listStatus.Clear()

	t := a.currentTamagotchi

	// Status header
	status := "🟢 Alive"
	if !t.IsAlive {
		status = "🔴 Dead"
	}
	listStatus.AddItem(fmt.Sprintf("Status: %s", status), "", 0, nil)

	// Basic info
	listStatus.AddItem(fmt.Sprintf("Name: %s", t.Name), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Age: %d days", t.Age), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Stage: %s", t.Stage), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Weight: %.1f grams", t.Weight), "", 0, nil)

	// Stats with visual bars
	listStatus.AddItem("", "", 0, nil) // Empty line
	listStatus.AddItem("=== STATS ===", "", 0, nil)

	// Hunger bar
	hungerBar := a.createProgressBar(t.Hunger, 100)
	hungerColor := "🟢"
	if t.Hunger > 80 {
		hungerColor = "🔴"
	} else if t.Hunger > 60 {
		hungerColor = "🟡"
	}
	listStatus.AddItem(fmt.Sprintf("Hunger: %s %s", hungerColor, hungerBar), "", 0, nil)

	// Happiness bar
	happinessBar := a.createProgressBar(t.Happiness, 100)
	happinessColor := "😊"
	if t.Happiness < 20 {
		happinessColor = "😢"
	} else if t.Happiness < 50 {
		happinessColor = "😐"
	}
	listStatus.AddItem(fmt.Sprintf("Happiness: %s %s", happinessColor, happinessBar), "", 0, nil)

	// Health bar
	healthBar := a.createProgressBar(t.Health, 100)
	healthColor := "🟢"
	if t.Health < 30 {
		healthColor = "🔴"
	} else if t.Health < 70 {
		healthColor = "🟡"
	}
	listStatus.AddItem(fmt.Sprintf("Health: %s %s", healthColor, healthBar), "", 0, nil)

	// Energy bar
	energyBar := a.createProgressBar(t.Energy, 100)
	energyColor := "⚡"
	if t.Energy < 20 {
		energyColor = "😴"
	} else if t.Energy < 50 {
		energyColor = "😪"
	}
	listStatus.AddItem(fmt.Sprintf("Energy: %s %s", energyColor, energyBar), "", 0, nil)

	// Last actions
	listStatus.AddItem("", "", 0, nil) // Empty line
	listStatus.AddItem("=== LAST ACTIONS ===", "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Last Fed: %s", t.LastFed.Format("15:04")), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Last Play: %s", t.LastPlay.Format("15:04")), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Last Sleep: %s", t.LastSleep.Format("15:04")), "", 0, nil)

	// Created date
	listStatus.AddItem("", "", 0, nil) // Empty line
	listStatus.AddItem("=== INFO ===", "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Created: %s", t.Created.Format("2006-01-02 15:04")), "", 0, nil)
	listStatus.AddItem(fmt.Sprintf("Time Alive: %s", time.Since(t.Created).Round(time.Hour)), "", 0, nil)
}

func (a *App) createProgressBar(current, max int) string {
	const barLength = 20
	filled := (current * barLength) / max
	if filled > barLength {
		filled = barLength
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
	return fmt.Sprintf("[%d%%] %s", current, bar)
}

func (a *App) statusPage() (title string, content tview.Primitive) {
	listStatus := a.viewsList["status"]
	if listStatus == nil {
		listStatus = getList()
		a.viewsList["status"] = listStatus
	}

	a.generateStatusList(listStatus)

	title = statusSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listStatus, 0, 1, true), 0, 1, true)
}
