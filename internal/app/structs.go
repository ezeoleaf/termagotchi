package app

import "time"

type Tamagotchi struct {
	Name      string
	Age       int
	Hunger    int
	Happiness int
	Health    int
	Energy    int
	Weight    float64
	Stage     string
	Created   time.Time
	LastFed   time.Time
	LastPlay  time.Time
	LastSleep time.Time
	IsAlive   bool
}

type GameEvent struct {
	Type      string
	Message   string
	Timestamp time.Time
}

type Food struct {
	Name       string
	Nutrition  int
	Happiness  int
	Energy     int
	WeightGain float64
}
