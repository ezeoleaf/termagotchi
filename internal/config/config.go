package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	Tamagotchi TamagotchiConfig `yaml:"tamagotchi"`
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

	// Update current login time
	cfg.App.CurrentLogin = time.Now()

	return cfg, nil
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
