package app

import (
	"github.com/rivo/tview"
)

func (a *App) generateHelpList(listHelp *tview.List) {
	listHelp.Clear()

	listHelp.AddItem("=== TERMAGOTCHI HELP ===", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🎮 GAME OVERVIEW", "", 0, nil)
	listHelp.AddItem("Termagotchi is a terminal-based Tamagotchi simulation.", "", 0, nil)
	listHelp.AddItem("Take care of your digital pet by feeding, playing, and sleeping.", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("📊 STATS EXPLANATION", "", 0, nil)
	listHelp.AddItem("Hunger: 0 = Full, 100 = Starving", "", 0, nil)
	listHelp.AddItem("Happiness: 0 = Very Sad, 100 = Very Happy", "", 0, nil)
	listHelp.AddItem("Health: 0 = Sick, 100 = Healthy", "", 0, nil)
	listHelp.AddItem("Energy: 0 = Tired, 100 = Energetic", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🔄 STAGES OF LIFE", "", 0, nil)
	listHelp.AddItem("Egg → Baby → Child → Teen → Adult", "", 0, nil)
	listHelp.AddItem("Your tamagotchi evolves based on age.", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("⚡ GAME MECHANICS", "", 0, nil)
	listHelp.AddItem("• Stats change automatically over time", "", 0, nil)
	listHelp.AddItem("• Keep hunger low and happiness high", "", 0, nil)
	listHelp.AddItem("• Low health can lead to death", "", 0, nil)
	listHelp.AddItem("• Energy is needed for playing", "", 0, nil)
	listHelp.AddItem("• Configure game speed and overnight grace in config.yml", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🎯 TIPS FOR SUCCESS", "", 0, nil)
	listHelp.AddItem("• Feed regularly to prevent hunger", "", 0, nil)
	listHelp.AddItem("• Play games to increase happiness", "", 0, nil)
	listHelp.AddItem("• Put to sleep when energy is low", "", 0, nil)
	listHelp.AddItem("• Balance all stats for best health", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("⌨️  KEYBOARD CONTROLS", "", 0, nil)
	listHelp.AddItem("Ctrl+S: Status - View tamagotchi stats", "", 0, nil)
	listHelp.AddItem("Ctrl+F: Feed - Give food to tamagotchi", "", 0, nil)
	listHelp.AddItem("Ctrl+P: Play - Play games with tamagotchi", "", 0, nil)
	listHelp.AddItem("Ctrl+L: Sleep - Put tamagotchi to sleep", "", 0, nil)
	listHelp.AddItem("Ctrl+E: Events - View game history", "", 0, nil)
	listHelp.AddItem("Ctrl+H: Help - Show this help page", "", 0, nil)
	listHelp.AddItem("Ctrl+R: Restart - Reset tamagotchi to new egg", "", 0, nil)
	listHelp.AddItem("Ctrl+C: Quit - Exit the game", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🎮 NAVIGATION", "", 0, nil)
	listHelp.AddItem("• Use arrow keys to navigate lists", "", 0, nil)
	listHelp.AddItem("• Press Enter to select items", "", 0, nil)
	listHelp.AddItem("• Use Ctrl+key shortcuts for quick access", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("💾 SAVE SYSTEM", "", 0, nil)
	listHelp.AddItem("Your tamagotchi progress is automatically saved.", "", 0, nil)
	listHelp.AddItem("Data is stored in your config directory.", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🔄 RESTART FEATURE", "", 0, nil)
	listHelp.AddItem("• Use Ctrl+R to restart with a new tamagotchi", "", 0, nil)
	listHelp.AddItem("• Confirmation modal will ask for your approval", "", 0, nil)
	listHelp.AddItem("• All progress will be reset to a fresh egg", "", 0, nil)
	listHelp.AddItem("• Useful if your tamagotchi dies or you want a fresh start", "", 0, nil)
	listHelp.AddItem("", "", 0, nil) // Empty line

	listHelp.AddItem("🐛 TROUBLESHOOTING", "", 0, nil)
	listHelp.AddItem("• If your tamagotchi dies, use Ctrl+R to restart", "", 0, nil)
	listHelp.AddItem("• Check the Events page for recent activity", "", 0, nil)
	listHelp.AddItem("• Monitor stats regularly to prevent issues", "", 0, nil)
}

func (a *App) helpPage() (title string, content tview.Primitive) {
	listHelp := a.viewsList["help"]
	if listHelp == nil {
		listHelp = getList()
		a.viewsList["help"] = listHelp
	}

	a.generateHelpList(listHelp)

	title = helpSection
	return title, tview.NewFlex().
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(listHelp, 0, 1, true), 0, 1, true)
}
