package burgoking

import (
	"math/rand"
	"time"
)

var (
	RestaurantCodes = []int{
		22365, // Auchan Leers - Leers
		23911, // Auchan v2 - Villeneuve d'Ascq
		21109, // Euralille - Lille
		19974, // Gare Saint-Lazare - Paris
		24191, // Zone Commerciale Grand Tour 2 - Sainte-Eulalie
		22118, // Zone Commerciale de l'Épinette - Seclin
	}
)

type Meal interface {
	Restaurant() int
	Date() time.Time
}

type RandomMeal struct{}

func (rm *RandomMeal) Restaurant() int {
	return rand.Intn(len(RestaurantCodes))
}

func (rm *RandomMeal) Date() time.Time {
	return time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour).Add(11*time.Hour + time.Duration(rand.Int63n(int64(3*time.Hour))))
}
