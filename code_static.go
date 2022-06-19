package burgoking

import (
	"fmt"
	"time"
	"math/rand"
)

func GenerateCodeStatic(meal *Meal) (code string, err error) {
	var prefix string
	switch time.Now().Month() {
	case time.January:
		prefix = "BB"
	case time.February:
		prefix = "LS"
	case time.March:
		prefix = "JH"
	case time.April:
		prefix = "PL"
	case time.May:
		prefix = "BK"
	case time.June:
		prefix = "WH"
	case time.July:
		prefix = "FF"
	case time.August:
		prefix = "BF"
	case time.September:
		prefix = "CF"
	case time.October:
		prefix = "CK"
	case time.November:
		prefix = "CB"
	case time.December:
		prefix = "VM"
	}
	source = rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s%d", prefix, 10_000 + source.Intn(89_999)), nil
}