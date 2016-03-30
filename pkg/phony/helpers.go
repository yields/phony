package phony

import (
	"math/rand"
	"strconv"
	"time"
)

func stringsToFloats(args []string) []float64 {
	floatArgs := []float64{}

	for i := 0; i < len(args); i++ {
		num, err := strconv.ParseFloat(args[i], 64)

		if err != nil {
			break
		}

		floatArgs = append(floatArgs, num)
	}

	return floatArgs
}

func CreateRand() *rand.Rand {
	source := rand.NewSource(time.Now().UnixNano())
	return rand.New(source)
}
