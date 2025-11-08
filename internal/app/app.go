package app

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ezeoleaf/termagotchi/internal/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const defaultUpdateInterval = 30 * time.Second

// App contains the tview application, the layout for the display, and the loaded config
type App struct {
	TApp              *tview.Application
	TLayout           *tview.Flex
	Config            *config.Config
	viewsList         map[string]*tview.List
	spriteView        *tview.TextView
	currentTamagotchi *Tamagotchi
	gameEvents        []GameEvent
	modal             *tview.Modal

	stateMu         sync.RWMutex
	eventsMu        sync.Mutex
	uiMu            sync.Mutex
	tuiRunning      bool
	uiReady         bool
	uiReadyOnce     sync.Once
	needsRefresh    bool
	timeAccumulator time.Duration
	updateInterval  time.Duration
}

// NewApp returns an instance of the application, initialized with the provided config
func NewApp(cfg *config.Config) *App {
	app := &App{
		TApp:           tview.NewApplication(),
		Config:         cfg,
		viewsList:      make(map[string]*tview.List),
		gameEvents:     make([]GameEvent, 0),
		updateInterval: defaultUpdateInterval,
	}

	app.initializeStateFromConfig()

	pages, info := app.getPagesInfo()

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	app.TApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
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

func (a *App) initializeStateFromConfig() {
	if a.Config.App.LastLogin.IsZero() {
		a.currentTamagotchi = newDefaultTamagotchi(randomName())
		a.updateConfigFromState()
	} else {
		cfg := a.Config.Tamagotchi
		a.currentTamagotchi = &Tamagotchi{
			Name:      cfg.Name,
			Age:       cfg.Age,
			Hunger:    cfg.Hunger,
			Happiness: cfg.Happiness,
			Health:    cfg.Health,
			Energy:    cfg.Energy,
			Weight:    cfg.Weight,
			Stage:     cfg.Stage,
			Created:   cfg.Created,
			LastFed:   cfg.LastFed,
			LastPlay:  cfg.LastPlay,
			LastSleep: cfg.LastSleep,
			IsAlive:   cfg.IsAlive,
		}
	}

	a.applyOfflineProgress()
}

func (a *App) Run() {
	a.uiMu.Lock()
	a.tuiRunning = true
	a.uiMu.Unlock()

	a.TApp.SetAfterDrawFunc(func(screen tcell.Screen) {
		a.markUIReady()
		a.TApp.SetAfterDrawFunc(nil)
	})

	app := a.TApp.SetRoot(a.TLayout, true).EnableMouse(true)

	defer func() {
		a.uiMu.Lock()
		a.tuiRunning = false
		a.uiMu.Unlock()

		now := time.Now()
		a.Config.App.LastLogin = now
		a.Config.App.CurrentLogin = now

		a.updateConfigFromState()

		if err := config.SaveConfig(a.Config); err != nil {
			log.Printf("failed to save config: %v", err)
		}
	}()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (a *App) showRestartModal() {
	if a.modal != nil {
		return // Modal already showing
	}

	modal := tview.NewModal().
		SetText("Are you sure you want to restart?\n\nThis will reset your tamagotchi to a new egg.\nAll progress will be lost!").
		AddButtons([]string{"Cancel", "Restart"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Restart" {
				a.restartTamagotchi()
			}
			a.TApp.SetRoot(a.TLayout, true).SetFocus(a.TLayout)
			a.modal = nil
		})

	a.modal = modal
	a.TApp.SetRoot(modal, true).SetFocus(modal)
}

func (a *App) restartTamagotchi() {
	newName := randomName()
	newTamagotchi := newDefaultTamagotchi(newName)

	a.stateMu.Lock()
	a.currentTamagotchi = newTamagotchi
	a.timeAccumulator = 0
	a.stateMu.Unlock()

	a.eventsMu.Lock()
	a.gameEvents = make([]GameEvent, 0)
	a.eventsMu.Unlock()

	a.addGameEvent("RESTART", fmt.Sprintf("Started a new tamagotchi named %s! 🥚", newName))
	a.updateConfigFromState()
}

func (a *App) gameLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastTick := time.Now()
	for now := range ticker.C {
		elapsed := now.Sub(lastTick)
		lastTick = now
		if elapsed < 0 {
			elapsed = 0
		}
		a.advanceTime(elapsed)
		a.refreshUI()
	}
}

func (a *App) advanceTime(elapsed time.Duration) bool {
	if elapsed <= 0 {
		return false
	}

	var ticks int

	a.stateMu.Lock()
	if a.currentTamagotchi != nil {
		a.timeAccumulator += elapsed
		ticks = int(a.timeAccumulator / a.updateInterval)
		if ticks > 0 {
			a.timeAccumulator -= time.Duration(ticks) * a.updateInterval
			for i := 0; i < ticks; i++ {
				a.applyTickLocked()
				if a.currentTamagotchi == nil || !a.currentTamagotchi.IsAlive {
					break
				}
			}
		}
	}
	a.stateMu.Unlock()

	if ticks > 0 {
		a.requestRefresh()
		return true
	}

	return false
}

func (a *App) applyTickLocked() {
	if a.currentTamagotchi == nil || !a.currentTamagotchi.IsAlive {
		return
	}

	t := a.currentTamagotchi

	t.Hunger = min(100, t.Hunger+5)

	if t.Hunger > 80 {
		t.Happiness = max(0, t.Happiness-2)
	}

	t.Energy = max(0, t.Energy-3)

	if t.Hunger > 90 || t.Happiness < 10 {
		t.Health = max(0, t.Health-1)
	}

	if t.Health <= 0 {
		if t.IsAlive {
			t.IsAlive = false
			a.addGameEvent("DEATH", "Your tamagotchi has passed away... 💔")
		}
		return
	}

	ageInHours := int(time.Since(t.Created).Hours())
	if ageInHours < 0 {
		ageInHours = 0
	}
	t.Age = ageInHours / 24

	oldStage := t.Stage
	a.updateStageLocked(oldStage)
}

func (a *App) updateStageLocked(previousStage string) {
	if a.currentTamagotchi == nil {
		return
	}

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

	if previousStage != a.currentTamagotchi.Stage {
		a.addGameEvent("EVOLUTION", fmt.Sprintf("Your tamagotchi evolved to %s! 🎉", a.currentTamagotchi.Stage))
	}
}

func (a *App) applyOfflineProgress() {
	if a.Config.App.LastLogin.IsZero() {
		return
	}

	elapsed := a.Config.App.CurrentLogin.Sub(a.Config.App.LastLogin)
	if elapsed <= 0 {
		return
	}

	if a.advanceTime(elapsed) {
		a.addGameEvent("PROGRESS", fmt.Sprintf("Time passed while you were away: %s.", formatDuration(elapsed)))
		a.updateConfigFromState()
	}
}

func (a *App) requestRefresh() {
	a.uiMu.Lock()
	a.needsRefresh = true
	ready := a.uiReady
	running := a.tuiRunning
	a.uiMu.Unlock()

	if ready && running {
		go a.refreshUI()
	}
}

func (a *App) refreshUI() {
	a.uiMu.Lock()
	if !a.tuiRunning || !a.uiReady || !a.needsRefresh || a.TApp == nil {
		a.uiMu.Unlock()
		return
	}
	a.needsRefresh = false
	a.uiMu.Unlock()

	a.TApp.QueueUpdateDraw(func() {
		if list := a.viewsList["status"]; list != nil {
			a.generateStatusList(list)
		}
		if a.spriteView != nil {
			a.updateSpriteView(a.spriteView)
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

func (a *App) updateConfigFromState() {
	a.stateMu.RLock()
	if a.currentTamagotchi == nil {
		a.stateMu.RUnlock()
		return
	}
	t := *a.currentTamagotchi
	a.stateMu.RUnlock()

	a.Config.Tamagotchi.Name = t.Name
	a.Config.Tamagotchi.Age = t.Age
	a.Config.Tamagotchi.Hunger = t.Hunger
	a.Config.Tamagotchi.Happiness = t.Happiness
	a.Config.Tamagotchi.Health = t.Health
	a.Config.Tamagotchi.Energy = t.Energy
	a.Config.Tamagotchi.Weight = t.Weight
	a.Config.Tamagotchi.Stage = t.Stage
	a.Config.Tamagotchi.Created = t.Created
	a.Config.Tamagotchi.LastFed = t.LastFed
	a.Config.Tamagotchi.LastPlay = t.LastPlay
	a.Config.Tamagotchi.LastSleep = t.LastSleep
	a.Config.Tamagotchi.IsAlive = t.IsAlive
}

func (a *App) tamagotchiSnapshot() (Tamagotchi, bool) {
	a.stateMu.RLock()
	defer a.stateMu.RUnlock()

	if a.currentTamagotchi == nil {
		return Tamagotchi{}, false
	}

	return *a.currentTamagotchi, true
}

func (a *App) eventsSnapshot() []GameEvent {
	a.eventsMu.Lock()
	defer a.eventsMu.Unlock()

	if len(a.gameEvents) == 0 {
		return nil
	}

	events := make([]GameEvent, len(a.gameEvents))
	copy(events, a.gameEvents)
	return events
}

func newDefaultTamagotchi(name string) *Tamagotchi {
	now := time.Now()
	return &Tamagotchi{
		Name:      name,
		Age:       0,
		Hunger:    50,
		Happiness: 50,
		Health:    100,
		Energy:    100,
		Weight:    50.0,
		Stage:     "egg",
		Created:   now,
		LastFed:   now,
		LastPlay:  now,
		LastSleep: now,
		IsAlive:   true,
	}
}

func formatDuration(d time.Duration) string {
	if d >= 24*time.Hour {
		days := d / (24 * time.Hour)
		hours := (d % (24 * time.Hour)) / time.Hour
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if d >= time.Hour {
		hours := d / time.Hour
		minutes := (d % time.Hour) / time.Minute
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	if d >= time.Minute {
		minutes := d / time.Minute
		seconds := (d % time.Minute) / time.Second
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	seconds := d / time.Second
	if seconds <= 0 {
		seconds = 1
	}
	return fmt.Sprintf("%ds", seconds)
}

func (a *App) addGameEvent(eventType, message string) {
	event := GameEvent{
		Type:      eventType,
		Message:   message,
		Timestamp: time.Now(),
	}

	a.eventsMu.Lock()
	a.gameEvents = append(a.gameEvents, event)

	if len(a.gameEvents) > 50 {
		a.gameEvents = a.gameEvents[len(a.gameEvents)-50:]
	}
	a.eventsMu.Unlock()

	a.requestRefresh()
}

func (a *App) markUIReady() {
	a.uiReadyOnce.Do(func() {
		a.uiMu.Lock()
		a.uiReady = true
		needs := a.needsRefresh
		a.uiMu.Unlock()

		if needs {
			go a.refreshUI()
		}
	})
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
