# Termagotchi 🐾

A terminal-based Tamagotchi simulation written in Go, featuring a beautiful TUI interface.

<img width="1188" height="725" alt="Screenshot from 2025-08-10 13-41-38" src="https://github.com/user-attachments/assets/f92eae52-0a2e-4816-ac62-e047cfae801d" />


## Features

- 🎮 **Interactive TUI**: Beautiful terminal user interface with keyboard navigation
- 🐾 **Digital Pet Care**: Feed, play, and put your tamagotchi to sleep
- 📊 **Real-time Stats**: Monitor hunger, happiness, health, and energy levels
- 🔄 **Life Stages**: Watch your tamagotchi evolve from egg to adult
- 📝 **Event History**: Track all interactions and milestones
- 💾 **Auto-save**: Progress is automatically saved to your config directory
- ⏰ **Time-based Mechanics**: Stats change over time, requiring regular care
- 🔄 **Restart Feature**: Reset to a new tamagotchi with confirmation modal

## Installation

### Homebrew (recommended, after release)

```sh
brew tap ezeoleaf/tap
brew install termagotchi
```

### Download a Release

- Go to [Releases](https://github.com/ezeoleaf/termagotchi/releases) and download the binary for your OS.
- Unpack and move it to a directory in your `$PATH` (e.g., `/usr/local/bin`).

### Build from Source

#### Prerequisites

- Go 1.24.0 or later

```bash
git clone <repository-url>
cd termagotchi
go mod tidy
go build -o termagotchi cmd/termagotchi/main.go
./termagotchi
```

## Usage

### Controls

- **Ctrl+S**: Status - View tamagotchi stats
- **Ctrl+F**: Feed - Give food to tamagotchi
- **Ctrl+P**: Play - Play games with tamagotchi
- **Ctrl+L**: Sleep - Put tamagotchi to sleep
- **Ctrl+E**: Events - View game history
- **Ctrl+H**: Help - Show help page
- **Ctrl+R**: Restart - Reset tamagotchi to new egg
- **Ctrl+C**: Quit - Exit the game

### Navigation

- Use arrow keys to navigate lists
- Press Enter to select items
- Use Ctrl+key shortcuts for quick access

## Game Mechanics

### Stats

- **Hunger**: 0 = Full, 100 = Starving
- **Happiness**: 0 = Very Sad, 100 = Very Happy
- **Health**: 0 = Sick, 100 = Healthy
- **Energy**: 0 = Tired, 100 = Energetic

### Life Stages

1. **Egg** (0-1 days)
2. **Baby** (1-3 days)
3. **Child** (3-7 days)
4. **Teen** (7-14 days)
5. **Adult** (14+ days)

### Food Types

- 🍎 Apple: Good nutrition, low weight gain
- 🍕 Pizza: High nutrition and happiness
- 🥗 Salad: Healthy option
- 🍔 Burger: High nutrition but heavy
- 🍦 Ice Cream: High happiness boost
- 🥕 Carrot: Balanced nutrition
- 🍫 Chocolate: Happiness boost
- 🥩 Steak: Maximum nutrition

### Games

- 🎾 Play Ball: Classic fun
- 🏃‍♂️ Run Around: Good exercise
- 🎵 Sing Songs: Low energy, high happiness
- 🎨 Draw Pictures: Creative fun
- 🧩 Solve Puzzle: Mental stimulation
- 🎭 Dance Party: High energy fun
- 📚 Read Books: Educational
- 🎪 Play Hide & Seek: Interactive fun

### Sleep Options

- 😴 Short Nap (30 min): Quick energy boost
- 😪 Medium Sleep (2 hours): Balanced rest
- 😴 Long Sleep (6 hours): Good recovery
- 😴 Full Night (8 hours): Complete restoration

## Tips for Success

1. **Feed Regularly**: Keep hunger below 80 to maintain happiness
2. **Play Often**: Increase happiness and health through games
3. **Rest When Needed**: Put to sleep when energy is low
4. **Monitor Health**: Low health can lead to death
5. **Balance Stats**: Keep all stats in good ranges
6. **Restart When Needed**: Use Ctrl+R to start fresh if your tamagotchi dies

## Configuration

The game automatically saves your progress to:

- **macOS**: `~/Library/Application Support/termagotchi/config.yml`
- **Linux**: `~/.config/termagotchi/config.yml`
- **Windows**: `%APPDATA%\termagotchi\config.yml`

## Development

### Project Structure

```
termagotchi/
├── cmd/
│   └── termagotchi/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   ├── structs.go
│   │   ├── pages.go
│   │   ├── status.go
│   │   ├── feed.go
│   │   ├── play.go
│   │   ├── sleep.go
│   │   ├── events.go
│   │   └── help.go
│   └── config/
│       └── config.go
├── go.mod
├── go.sum
└── README.md
```

### Dependencies

- `github.com/gdamore/tcell/v2`: Terminal UI framework
- `github.com/rivo/tview`: TUI components
- `gopkg.in/yaml.v3`: Configuration file handling

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Inspired by the classic Tamagotchi digital pet
- Built with the amazing `tview` TUI framework
- Thanks to the Go community for excellent libraries
