package app

import (
	"testing"
	"time"
)

func TestUpdateIntervalForSpeed(t *testing.T) {
	tests := []struct {
		speed float64
		want  time.Duration
	}{
		{speed: 0.5, want: 60 * time.Second},
		{speed: 1, want: 30 * time.Second},
		{speed: 2, want: 15 * time.Second},
	}

	for _, tt := range tests {
		got, err := updateIntervalForSpeed(tt.speed)
		if err != nil {
			t.Fatalf("update interval for speed %v: %v", tt.speed, err)
		}
		if got != tt.want {
			t.Errorf("update interval for speed %v = %s, want %s", tt.speed, got, tt.want)
		}
	}

	for _, speed := range []float64{0, -1} {
		if _, err := updateIntervalForSpeed(speed); err == nil {
			t.Errorf("expected speed %v to be rejected", speed)
		}
	}
}

func TestOfflineProgressDuration(t *testing.T) {
	location := time.FixedZone("local", 0)
	tests := []struct {
		name        string
		logout      time.Time
		login       time.Time
		grace       time.Duration
		wantElapsed time.Duration
		wantForgive time.Duration
	}{
		{
			name:        "same day does not receive grace",
			logout:      time.Date(2026, 10, 28, 18, 0, 0, 0, location),
			login:       time.Date(2026, 10, 28, 20, 0, 0, 0, location),
			grace:       8 * time.Hour,
			wantElapsed: 2 * time.Hour,
		},
		{
			name:        "next day forgives configured duration",
			logout:      time.Date(2026, 10, 28, 23, 0, 0, 0, location),
			login:       time.Date(2026, 10, 29, 8, 0, 0, 0, location),
			grace:       8 * time.Hour,
			wantElapsed: time.Hour,
			wantForgive: 8 * time.Hour,
		},
		{
			name:        "grace cannot exceed time away",
			logout:      time.Date(2026, 10, 28, 23, 0, 0, 0, location),
			login:       time.Date(2026, 10, 29, 1, 0, 0, 0, location),
			grace:       8 * time.Hour,
			wantForgive: 2 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elapsed, forgiven := offlineProgressDuration(tt.logout, tt.login, tt.grace)
			if elapsed != tt.wantElapsed || forgiven != tt.wantForgive {
				t.Errorf("got elapsed %s, forgiven %s; want elapsed %s, forgiven %s",
					elapsed, forgiven, tt.wantElapsed, tt.wantForgive)
			}
		})
	}
}
