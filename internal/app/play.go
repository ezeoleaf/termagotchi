package app

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

type Game struct {
	Name       string
	Happiness  int
	Energy     int
	Health     int
	WeightLoss float64
}

var availableGames = []Game{
	{Name: "🎾 Play Ball", Happiness: 20, Energy: -15, Health: 5, WeightLoss: 0.5},
	{Name: "🏃‍♂️ Run Around", Happiness: 15, Energy: -25, Health: 10, WeightLoss: 1.0},
	{Name: "🎵 Sing Songs", Happiness: 25, Energy: -5, Health: 3, WeightLoss: 0.1},
	{Name: "🎨 Draw Pictures", Happiness: 30, Energy: -10, Health: 2, WeightLoss: 0.2},
	{Name: "🧩 Solve Puzzle", Happiness: 35, Energy: -20, Health: 8, WeightLoss: 0.3},
	{Name: "🎭 Dance Party", Happiness: 40, Energy: -30, Health: 12, WeightLoss: 1.5},
	{Name: "📚 Read Books", Happiness: 15, Energy: -5, Health: 5, WeightLoss: 0.1},
	{Name: "🎪 Play Hide & Seek", Happiness: 25, Energy: -20, Health: 7, WeightLoss: 0.8},
}

func (a *App) generatePlayList(listPlay *tview.List) {
	listPlay.Clear()

	t, ok := a.tamagotchiSnapshot()
	if !ok {
		listPlay.AddItem("No tamagotchi available.", "", 0, nil)
		return
	}

	if !t.IsAlive {
		listPlay.AddItem("Your tamagotchi has passed away... 💔", "", 0, nil)
		listPlay.AddItem("Cannot play with a dead tamagotchi", "", 0, nil)
		return
	}

	// Check if tamagotchi has enough energy to play
	if t.Energy < 10 {
		listPlay.AddItem("😴 Your tamagotchi is too tired to play!", "", 0, nil)
		listPlay.AddItem("Try putting it to sleep first (Ctrl+L)", "", 0, nil)
		listPlay.AddItem("", "", 0, nil) // Empty line
	}

	listPlay.AddItem("=== AVAILABLE GAMES ===", "", 0, nil)
	listPlay.AddItem("", "", 0, nil) // Empty line

	for i, game := range availableGames {
		gameIndex := i // Capture the index for the closure
		energyChange := ""
		if game.Energy < 0 {
			energyChange = fmt.Sprintf("Energy: %d", game.Energy)
		} else {
			energyChange = fmt.Sprintf("Energy: +%d", game.Energy)
		}

		listPlay.AddItem(
			fmt.Sprintf("%s (Happiness: +%d, %s, Health: +%d, Weight: -%.1fg)",
				game.Name, game.Happiness, energyChange, game.Health, game.WeightLoss),
			"",
			0,
			func() { a.playWithTamagotchi(gameIndex) },
		)
	}

	listPlay.AddItem("", "", 0, nil) // Empty line
	listPlay.AddItem("=== PLAYING INFO ===", "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Happiness: %d/100", t.Happiness), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Energy: %d/100", t.Energy), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Weight: %.1f grams", t.Weight), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Last Play: %s", t.LastPlay.Format("15:04")), "", 0, nil)
}

func (a *App) playWithTamagotchi(gameIndex int) {
	if gameIndex < 0 || gameIndex >= len(availableGames) {
		return
	}

	game := availableGames[gameIndex]
	now := time.Now()

	a.stateMu.Lock()
	if a.currentTamagotchi == nil || !a.currentTamagotchi.IsAlive {
		a.stateMu.Unlock()
		return
	}

	if a.currentTamagotchi.Energy < 10 {
		a.stateMu.Unlock()
		return
	}

	a.currentTamagotchi.Happiness = min(100, a.currentTamagotchi.Happiness+game.Happiness)
	a.currentTamagotchi.Energy = max(0, min(100, a.currentTamagotchi.Energy+game.Energy))
	a.currentTamagotchi.Health = min(100, a.currentTamagotchi.Health+game.Health)
	if a.currentTamagotchi.Weight-game.WeightLoss < 10.0 {
		a.currentTamagotchi.Weight = 10.0
	} else {
		a.currentTamagotchi.Weight -= game.WeightLoss
	}
	a.currentTamagotchi.LastPlay = now
	a.stateMu.Unlock()

	a.updateConfigFromState()
	a.addGameEvent("PLAY", fmt.Sprintf("Played %s! Happiness +%d, Energy %d", game.Name, game.Happiness, game.Energy))
}

func (a *App) playPage() (title string, content tview.Primitive) {
	listPlay := a.viewsList["play"]
	if listPlay == nil {
		listPlay = getList()
		a.viewsList["play"] = listPlay
	}

	a.generatePlayList(listPlay)

	title = playSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listPlay, 0, 1, true), 0, 1, true)
}
