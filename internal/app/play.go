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

	if !a.currentTamagotchi.IsAlive {
		listPlay.AddItem("Your tamagotchi has passed away... 💔", "", 0, nil)
		listPlay.AddItem("Cannot play with a dead tamagotchi", "", 0, nil)
		return
	}

	// Check if tamagotchi has enough energy to play
	if a.currentTamagotchi.Energy < 10 {
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
			func() {
				a.playWithTamagotchi(gameIndex)
				a.generatePlayList(listPlay)
			},
		)
	}

	listPlay.AddItem("", "", 0, nil) // Empty line
	listPlay.AddItem("=== PLAYING INFO ===", "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Happiness: %d/100", a.currentTamagotchi.Happiness), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Energy: %d/100", a.currentTamagotchi.Energy), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Current Weight: %.1f grams", a.currentTamagotchi.Weight), "", 0, nil)
	listPlay.AddItem(fmt.Sprintf("Last Play: %s", a.currentTamagotchi.LastPlay.Format("15:04")), "", 0, nil)
}

func (a *App) playWithTamagotchi(gameIndex int) {
	if !a.currentTamagotchi.IsAlive {
		return
	}

	game := availableGames[gameIndex]

	// Check if tamagotchi has enough energy
	if a.currentTamagotchi.Energy < 10 {
		return
	}

	// Increase happiness
	a.currentTamagotchi.Happiness = min(100, a.currentTamagotchi.Happiness+game.Happiness)

	// Change energy (can be negative)
	a.currentTamagotchi.Energy = max(0, min(100, a.currentTamagotchi.Energy+game.Energy))

	// Increase health
	a.currentTamagotchi.Health = min(100, a.currentTamagotchi.Health+game.Health)

	// Decrease weight
	if a.currentTamagotchi.Weight-game.WeightLoss < 10.0 {
		a.currentTamagotchi.Weight = 10.0
	} else {
		a.currentTamagotchi.Weight -= game.WeightLoss
	}

	// Update last play time
	a.currentTamagotchi.LastPlay = time.Now()

	// Add game event
	a.addGameEvent("PLAY", fmt.Sprintf("Played %s! Happiness +%d, Energy %d", game.Name, game.Happiness, game.Energy))

	// Update status page if it exists
	if list := a.viewsList["status"]; list != nil {
		a.generateStatusList(list)
	}
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
