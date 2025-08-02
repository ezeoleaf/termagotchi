package app

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

var availableFoods = []Food{
	{Name: "🍎 Apple", Nutrition: 20, Happiness: 5, Energy: 10, WeightGain: 0.5},
	{Name: "🍕 Pizza", Nutrition: 40, Happiness: 15, Energy: 20, WeightGain: 2.0},
	{Name: "🥗 Salad", Nutrition: 15, Happiness: 3, Energy: 5, WeightGain: 0.2},
	{Name: "🍔 Burger", Nutrition: 50, Happiness: 20, Energy: 25, WeightGain: 3.0},
	{Name: "🍦 Ice Cream", Nutrition: 10, Happiness: 25, Energy: 15, WeightGain: 1.5},
	{Name: "🥕 Carrot", Nutrition: 25, Happiness: 8, Energy: 12, WeightGain: 0.3},
	{Name: "🍫 Chocolate", Nutrition: 15, Happiness: 30, Energy: 20, WeightGain: 1.0},
	{Name: "🥩 Steak", Nutrition: 60, Happiness: 10, Energy: 30, WeightGain: 4.0},
}

func (a *App) generateFeedList(listFeed *tview.List) {
	listFeed.Clear()

	if !a.currentTamagotchi.IsAlive {
		listFeed.AddItem("Your tamagotchi has passed away... 💔", "", 0, nil)
		listFeed.AddItem("Cannot feed a dead tamagotchi", "", 0, nil)
		return
	}

	listFeed.AddItem("=== AVAILABLE FOOD ===", "", 0, nil)
	listFeed.AddItem("", "", 0, nil) // Empty line

	for i, food := range availableFoods {
		foodIndex := i // Capture the index for the closure
		listFeed.AddItem(
			fmt.Sprintf("%s (Nutrition: %d, Happiness: %d, Energy: %d, Weight: +%.1fg)",
				food.Name, food.Nutrition, food.Happiness, food.Energy, food.WeightGain),
			"",
			0,
			func() {
				a.feedTamagotchi(foodIndex)
				a.generateFeedList(listFeed)
			},
		)
	}

	listFeed.AddItem("", "", 0, nil) // Empty line
	listFeed.AddItem("=== FEEDING INFO ===", "", 0, nil)
	listFeed.AddItem(fmt.Sprintf("Current Hunger: %d/100", a.currentTamagotchi.Hunger), "", 0, nil)
	listFeed.AddItem(fmt.Sprintf("Current Weight: %.1f grams", a.currentTamagotchi.Weight), "", 0, nil)
	listFeed.AddItem(fmt.Sprintf("Last Fed: %s", a.currentTamagotchi.LastFed.Format("15:04")), "", 0, nil)
}

func (a *App) feedTamagotchi(foodIndex int) {
	if !a.currentTamagotchi.IsAlive {
		return
	}

	food := availableFoods[foodIndex]

	// Reduce hunger
	a.currentTamagotchi.Hunger = max(0, a.currentTamagotchi.Hunger-food.Nutrition)

	// Increase happiness
	a.currentTamagotchi.Happiness = min(100, a.currentTamagotchi.Happiness+food.Happiness)

	// Increase energy
	a.currentTamagotchi.Energy = min(100, a.currentTamagotchi.Energy+food.Energy)

	// Increase weight
	a.currentTamagotchi.Weight += food.WeightGain

	// Update last fed time
	a.currentTamagotchi.LastFed = time.Now()

	// Add game event
	a.addGameEvent("FEED", fmt.Sprintf("Fed %s! Hunger -%d, Happiness +%d", food.Name, food.Nutrition, food.Happiness))

	// Update status page if it exists
	if list := a.viewsList["status"]; list != nil {
		a.generateStatusList(list)
	}
}

func (a *App) feedPage() (title string, content tview.Primitive) {
	listFeed := a.viewsList["feed"]
	if listFeed == nil {
		listFeed = getList()
		a.viewsList["feed"] = listFeed
	}

	a.generateFeedList(listFeed)

	title = feedSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listFeed, 0, 1, true), 0, 1, true)
}
