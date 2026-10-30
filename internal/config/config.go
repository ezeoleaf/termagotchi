package config

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	Game       GameConfig       `yaml:"game"`
	Tamagotchi TamagotchiConfig `yaml:"tamagotchi"`
}

type GameConfig struct {
	Speed       float64  `yaml:"speed"`
	GracePeriod Duration `yaml:"grace_period"`
}

type Duration time.Duration

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalYAML() (any, error) {
	return time.Duration(d).String(), nil
}

func (d *Duration) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.ScalarNode {
		return fmt.Errorf("grace_period must be a duration string")
	}

	duration, err := time.ParseDuration(node.Value)
	if err != nil {
		return fmt.Errorf("invalid grace_period %q: %w", node.Value, err)
	}
	if duration < 0 {
		return fmt.Errorf("grace_period cannot be negative")
	}

	*d = Duration(duration)
	return nil
}

type AppConfig struct {
	LastLogin     time.Time `yaml:"last_login"`
	CurrentLogin  time.Time `yaml:"current_login"`
	SaveDirectory string    `yaml:"save_directory"`
}

type TamagotchiConfig struct {
	Name      string    `yaml:"name"`
	Age       int       `yaml:"age"`
	Hunger    int       `yaml:"hunger"`    // 0-100, 0 = full, 100 = starving
	Happiness int       `yaml:"happiness"` // 0-100, 0 = very sad, 100 = very happy
	Health    int       `yaml:"health"`    // 0-100, 0 = sick, 100 = healthy
	Energy    int       `yaml:"energy"`    // 0-100, 0 = tired, 100 = energetic
	Weight    float64   `yaml:"weight"`    // in grams
	Stage     string    `yaml:"stage"`     // egg, baby, child, teen, adult
	Created   time.Time `yaml:"created"`
	LastFed   time.Time `yaml:"last_fed"`
	LastPlay  time.Time `yaml:"last_play"`
	LastSleep time.Time `yaml:"last_sleep"`
	IsAlive   bool      `yaml:"is_alive"`
}

func LoadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appConfigDir := filepath.Join(configDir, "termagotchi")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return nil, err
	}

	configPath := filepath.Join(appConfigDir, "config.yml")

	cfg := &Config{
		App: AppConfig{
			CurrentLogin:  time.Now(),
			SaveDirectory: appConfigDir,
		},
		Game: GameConfig{
			Speed:       1,
			GracePeriod: Duration(8 * time.Hour),
		},
		Tamagotchi: TamagotchiConfig{
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
		},
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Update current login time
	cfg.App.CurrentLogin = time.Now()

	return cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg == nil {
		return fmt.Errorf("config cannot be nil")
	}
	if cfg.Game.Speed <= 0 || math.IsNaN(cfg.Game.Speed) || math.IsInf(cfg.Game.Speed, 0) {
		return fmt.Errorf("game speed must be a finite number greater than zero")
	}
	if cfg.Game.GracePeriod.Duration() < 0 {
		return fmt.Errorf("game grace_period cannot be negative")
	}
	return nil
}

func SaveConfig(cfg *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appConfigDir := filepath.Join(configDir, "termagotchi")
	configPath := filepath.Join(appConfigDir, "config.yml")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
