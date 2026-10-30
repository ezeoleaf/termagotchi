package config

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestGameConfigYAMLRoundTrip(t *testing.T) {
	want := Config{
		Game: GameConfig{
			Speed:       1.5,
			GracePeriod: Duration(8 * time.Hour),
		},
	}

	data, err := yaml.Marshal(want)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}

	var got Config
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal config: %v", err)
	}

	if got.Game.Speed != want.Game.Speed {
		t.Errorf("speed = %v, want %v", got.Game.Speed, want.Game.Speed)
	}
	if got.Game.GracePeriod != want.Game.GracePeriod {
		t.Errorf("grace period = %v, want %v", got.Game.GracePeriod, want.Game.GracePeriod)
	}
}

func TestGracePeriodRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"-1h", "tomorrow"} {
		t.Run(value, func(t *testing.T) {
			var cfg Config
			err := yaml.Unmarshal([]byte("game:\n  grace_period: "+value+"\n"), &cfg)
			if err == nil {
				t.Fatalf("expected invalid grace period %q to fail", value)
			}
		})
	}
}

func TestGameConfigDefaultsWhenFieldsAreOmitted(t *testing.T) {
	cfg := Config{
		Game: GameConfig{
			Speed:       1,
			GracePeriod: Duration(8 * time.Hour),
		},
	}

	if err := yaml.Unmarshal([]byte("game:\n  speed: 2\n"), &cfg); err != nil {
		t.Fatalf("unmarshal partial game config: %v", err)
	}
	if cfg.Game.Speed != 2 {
		t.Errorf("speed = %v, want 2", cfg.Game.Speed)
	}
	if cfg.Game.GracePeriod != Duration(8*time.Hour) {
		t.Errorf("grace period = %v, want 8h default", cfg.Game.GracePeriod)
	}
}
