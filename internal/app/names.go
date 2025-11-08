package app

import (
	crand "crypto/rand"
	"math/big"
)

var tamagotchiNames = []string{
	"Leslie",
	"Ron",
	"Ann",
	"April",
	"Andy",
	"Ben",
	"Tom",
	"Donna",
	"Jerry",
	"Chris",
	"Craig",
	"Jean-Ralphio",
	"Mona-Lisa",
	"Perd",
	"Jeremy",
	"Bobby",
	"Tammy One",
	"Tammy Two",
	"Shauna Malwae-Tweep",
	"Joan",
}

func randomName() string {
	if len(tamagotchiNames) == 0 {
		return "Olaf"
	}
	upper := big.NewInt(int64(len(tamagotchiNames)))
	n, err := crand.Int(crand.Reader, upper)
	if err != nil {
		return tamagotchiNames[0]
	}
	return tamagotchiNames[n.Int64()]
}
