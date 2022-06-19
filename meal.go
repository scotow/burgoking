package burgoking

import (
	"math/rand"
	"time"
)

var (
	RestaurantCodes = []int{
		// 22365, // Auchan Leers - Leers
		// 23911, // Auchan v2 - Villeneuve d'Ascq
		// 21109, // Euralille - Lille
		// 19974, // Gare Saint-Lazare - Paris
		// 24191, // Zone Commerciale Grand Tour 2 - Sainte-Eulalie
		// 22118, // Zone Commerciale de l'Épinette - Seclin

		1677, // Tyler - Texas
	}

	source = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Meal struct {
	Restaurant int
	Date       time.Time
}

func RandomMeal() *Meal {
	return &Meal{
		RestaurantCodes[rand.Intn(len(RestaurantCodes))],
		time.Now().UTC().Add(-24 * time.Hour).Truncate(24 * time.Hour).Add(6*time.Hour + time.Duration(source.Int63n(int64(3*time.Hour)))),
	}
}
