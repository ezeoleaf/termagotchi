package app

import (
	"log"
	"time"

	"github.com/ezeoleaf/termagotchi/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// App contains the tview application, the layout for the display, and the loaded config
type App struct {
	TApp              *tview.Application
	TLayout           *tview.Flex
	Config            *config.Config
	viewsList         map[string]*tview.List
	currentTamagotchi *Tamagotchi
	gameEvents        []GameEvent
	modal             *tview.Modal
}

// NewApp returns an instance of the application, initialized with the provided config
func NewApp(cfg *config.Config) *App {
	app := &App{
		TApp:      tview.NewApplication(),
		Config:    cfg,
		viewsList: make(map[string]*tview.List),
		currentTamagotchi: &Tamagotchi{
			Name:      cfg.Tamagotchi.Name,
			Age:       cfg.Tamagotchi.Age,
			Hunger:    cfg.Tamagotchi.Hunger,
			Happiness: cfg.Tamagotchi.Happiness,
			Health:    cfg.Tamagotchi.Health,
			Energy:    cfg.Tamagotchi.Energy,
			Weight:    cfg.Tamagotchi.Weight,
			Stage:     cfg.Tamagotchi.Stage,
			Created:   cfg.Tamagotchi.Created,
			LastFed:   cfg.Tamagotchi.LastFed,
			LastPlay:  cfg.Tamagotchi.LastPlay,
			LastSleep: cfg.Tamagotchi.LastSleep,
			IsAlive:   cfg.Tamagotchi.IsAlive,
		},
		gameEvents: make([]GameEvent, 0),
	}

	pages, info := app.getPagesInfo()

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	app.TApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch ek := event.Key(); ek {
		case tcell.KeyCtrlH:
			app.goToSection(helpSection, info)
		case tcell.KeyCtrlS:
			app.goToSection(statusSection, info)
		case tcell.KeyCtrlF:
			app.goToSection(feedSection, info)
		case tcell.KeyCtrlP:
			app.goToSection(playSection, info)
		case tcell.KeyCtrlL:
			app.goToSection(sleepSection, info)
		case tcell.KeyCtrlE:
			app.goToSection(eventsSection, info)
		case tcell.KeyCtrlR:
			app.showRestartModal()
		}
		return event
	})

	app.TLayout = layout

	// Start the game loop
	go app.gameLoop()

	return app
}

func (a *App) Run() {
	if err := a.TApp.SetRoot(a.TLayout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	defer func() {
		a.Config.App.LastLogin = a.Config.App.CurrentLogin
		// Update config with current tamagotchi state
		a.Config.Tamagotchi.Name = a.currentTamagotchi.Name
		a.Config.Tamagotchi.Age = a.currentTamagotchi.Age
		a.Config.Tamagotchi.Hunger = a.currentTamagotchi.Hunger
		a.Config.Tamagotchi.Happiness = a.currentTamagotchi.Happiness
		a.Config.Tamagotchi.Health = a.currentTamagotchi.Health
		a.Config.Tamagotchi.Energy = a.currentTamagotchi.Energy
		a.Config.Tamagotchi.Weight = a.currentTamagotchi.Weight
		a.Config.Tamagotchi.Stage = a.currentTamagotchi.Stage
		a.Config.Tamagotchi.LastFed = a.currentTamagotchi.LastFed
		a.Config.Tamagotchi.LastPlay = a.currentTamagotchi.LastPlay
		a.Config.Tamagotchi.LastSleep = a.currentTamagotchi.LastSleep
		a.Config.Tamagotchi.IsAlive = a.currentTamagotchi.IsAlive

		err := config.SaveConfig(a.Config)
		if err != nil {
			log.Printf("failed to save config: %v", err)
		}
	}()
}

func (a *App) showRestartModal() {
	if a.modal != nil {
		return // Modal already showing
	}

	a.modal = tview.NewModal().
		SetText("Are you sure you want to restart?\n\nThis will reset your tamagotchi to a new egg.\nAll progress will be lost!").
		AddButtons([]string{"Cancel", "Restart"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Restart" {
				a.restartTamagotchi()
			}
			a.TApp.SetRoot(a.TLayout, true)
			a.modal = nil
		})

	a.TApp.SetRoot(a.modal, true)
}

func (a *App) restartTamagotchi() {
	// Create a new tamagotchi
	a.currentTamagotchi = &Tamagotchi{
		Name:      "Tammy",
		Age:       0,
		Hunger:    50,
		Happiness: 50,
		Health:    100,
		Energy:    100,
		Weight:    50.0,
		Stage:     "egg",
		Created:   time.Now(),
		LastFed:   time.Now(),
		LastPlay:  time.Now(),
		LastSleep: time.Now(),
		IsAlive:   true,
	}

	// Clear game events
	a.gameEvents = make([]GameEvent, 0)

	// Add restart event
	a.addGameEvent("RESTART", "Started a new tamagotchi! 🥚")

	// Update all views
	a.TApp.QueueUpdateDraw(func() {
		if list := a.viewsList["status"]; list != nil {
			a.generateStatusList(list)
		}
		if list := a.viewsList["feed"]; list != nil {
			a.generateFeedList(list)
		}
		if list := a.viewsList["play"]; list != nil {
			a.generatePlayList(list)
		}
		if list := a.viewsList["sleep"]; list != nil {
			a.generateSleepList(list)
		}
		if list := a.viewsList["events"]; list != nil {
			a.generateEventsList(list)
		}
	})
}

func (a *App) gameLoop() {
	ticker := time.NewTicker(30 * time.Second) // Update every 30 seconds
	defer ticker.Stop()

	for range ticker.C {
		a.updateTamagotchi()
	}
}

func (a *App) updateTamagotchi() {
	if !a.currentTamagotchi.IsAlive {
		return
	}

	// Increase hunger over time
	a.currentTamagotchi.Hunger = min(100, a.currentTamagotchi.Hunger+5)

	// Decrease happiness if hungry
	if a.currentTamagotchi.Hunger > 80 {
		a.currentTamagotchi.Happiness = max(0, a.currentTamagotchi.Happiness-2)
	}

	// Decrease energy over time
	a.currentTamagotchi.Energy = max(0, a.currentTamagotchi.Energy-3)

	// Decrease health if very hungry or very unhappy
	if a.currentTamagotchi.Hunger > 90 || a.currentTamagotchi.Happiness < 10 {
		a.currentTamagotchi.Health = max(0, a.currentTamagotchi.Health-1)
	}

	// Check if tamagotchi died
	if a.currentTamagotchi.Health <= 0 {
		a.currentTamagotchi.IsAlive = false
		a.addGameEvent("DEATH", "Your tamagotchi has passed away... 💔")
	}

	// Age calculation (1 day = 24 hours)
	ageInHours := int(time.Since(a.currentTamagotchi.Created).Hours())
	a.currentTamagotchi.Age = ageInHours / 24

	// Stage evolution
	a.updateStage()

	// Update UI
	a.TApp.QueueUpdateDraw(func() {
		if list := a.viewsList["status"]; list != nil {
			a.generateStatusList(list)
		}
	})
}

func (a *App) updateStage() {
	oldStage := a.currentTamagotchi.Stage

	switch {
	case a.currentTamagotchi.Age < 1:
		a.currentTamagotchi.Stage = "egg"
	case a.currentTamagotchi.Age < 3:
		a.currentTamagotchi.Stage = "baby"
	case a.currentTamagotchi.Age < 7:
		a.currentTamagotchi.Stage = "child"
	case a.currentTamagotchi.Age < 14:
		a.currentTamagotchi.Stage = "teen"
	default:
		a.currentTamagotchi.Stage = "adult"
	}

	if oldStage != a.currentTamagotchi.Stage {
		a.addGameEvent("EVOLUTION", "Your tamagotchi evolved to "+a.currentTamagotchi.Stage+"! 🎉")
	}
}

func (a *App) addGameEvent(eventType, message string) {
	event := GameEvent{
		Type:      eventType,
		Message:   message,
		Timestamp: time.Now(),
	}
	a.gameEvents = append(a.gameEvents, event)

	// Keep only last 50 events
	if len(a.gameEvents) > 50 {
		a.gameEvents = a.gameEvents[len(a.gameEvents)-50:]
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
