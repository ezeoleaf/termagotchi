package app

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

type SleepOption struct {
	Name       string
	Duration   time.Duration
	EnergyGain int
	HealthGain int
	Happiness  int
}

var sleepOptions = []SleepOption{
	{Name: "😴 Short Nap (30 min)", Duration: 30 * time.Minute, EnergyGain: 20, HealthGain: 5, Happiness: 5},
	{Name: "😪 Medium Sleep (2 hours)", Duration: 2 * time.Hour, EnergyGain: 50, HealthGain: 15, Happiness: 10},
	{Name: "😴 Long Sleep (6 hours)", Duration: 6 * time.Hour, EnergyGain: 80, HealthGain: 25, Happiness: 15},
	{Name: "😴 Full Night (8 hours)", Duration: 8 * time.Hour, EnergyGain: 100, HealthGain: 30, Happiness: 20},
}

func (a *App) generateSleepList(listSleep *tview.List) {
	listSleep.Clear()

	if !a.currentTamagotchi.IsAlive {
		listSleep.AddItem("Your tamagotchi has passed away... 💔", "", 0, nil)
		listSleep.AddItem("Cannot put a dead tamagotchi to sleep", "", 0, nil)
		return
	}

	listSleep.AddItem("=== SLEEP OPTIONS ===", "", 0, nil)
	listSleep.AddItem("", "", 0, nil) // Empty line

	for i, sleep := range sleepOptions {
		sleepIndex := i // Capture the index for the closure
		listSleep.AddItem(
			fmt.Sprintf("%s (Energy: +%d, Health: +%d, Happiness: +%d)",
				sleep.Name, sleep.EnergyGain, sleep.HealthGain, sleep.Happiness),
			"",
			0,
			func() {
				a.putTamagotchiToSleep(sleepIndex)
				a.generateSleepList(listSleep)
			},
		)
	}

	listSleep.AddItem("", "", 0, nil) // Empty line
	listSleep.AddItem("=== SLEEP INFO ===", "", 0, nil)
	listSleep.AddItem(fmt.Sprintf("Current Energy: %d/100", a.currentTamagotchi.Energy), "", 0, nil)
	listSleep.AddItem(fmt.Sprintf("Current Health: %d/100", a.currentTamagotchi.Health), "", 0, nil)
	listSleep.AddItem(fmt.Sprintf("Current Happiness: %d/100", a.currentTamagotchi.Happiness), "", 0, nil)
	listSleep.AddItem(fmt.Sprintf("Last Sleep: %s", a.currentTamagotchi.LastSleep.Format("15:04")), "", 0, nil)

	// Show sleep recommendation
	if a.currentTamagotchi.Energy < 30 {
		listSleep.AddItem("", "", 0, nil) // Empty line
		listSleep.AddItem("💡 Recommendation: Your tamagotchi is tired!", "", 0, nil)
		listSleep.AddItem("   Consider a longer sleep to restore energy.", "", 0, nil)
	}
}

func (a *App) putTamagotchiToSleep(sleepIndex int) {
	if !a.currentTamagotchi.IsAlive {
		return
	}

	sleep := sleepOptions[sleepIndex]

	// Increase energy
	a.currentTamagotchi.Energy = min(100, a.currentTamagotchi.Energy+sleep.EnergyGain)

	// Increase health
	a.currentTamagotchi.Health = min(100, a.currentTamagotchi.Health+sleep.HealthGain)

	// Increase happiness
	a.currentTamagotchi.Happiness = min(100, a.currentTamagotchi.Happiness+sleep.Happiness)

	// Update last sleep time
	a.currentTamagotchi.LastSleep = time.Now()

	// Add game event
	a.addGameEvent("SLEEP", fmt.Sprintf("Slept for %s! Energy +%d, Health +%d", sleep.Name, sleep.EnergyGain, sleep.HealthGain))

	// Update status page if it exists
	if list := a.viewsList["status"]; list != nil {
		a.generateStatusList(list)
	}
}

func (a *App) sleepPage() (title string, content tview.Primitive) {
	listSleep := a.viewsList["sleep"]
	if listSleep == nil {
		listSleep = getList()
		a.viewsList["sleep"] = listSleep
	}

	a.generateSleepList(listSleep)

	title = sleepSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listSleep, 0, 1, true), 0, 1, true)
}
