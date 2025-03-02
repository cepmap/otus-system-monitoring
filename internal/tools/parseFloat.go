package float

import (
	"strconv"
	"strings"
)

func ParseFloat(input string) float64 {
	buff := strings.ReplaceAll(input, ",", ".")
	output, err := strconv.ParseFloat(buff, 64)
	if err != nil {
		return 0.0
	}
	return output
}
