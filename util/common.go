package util

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

func Uniqid() string {
	now := time.Now()
	return fmt.Sprintf("%08x%08x", now.Unix(), now.UnixNano()%0x100000)
}

func CompareFloat(floatA, floatB float64, place int32) int {
	deciA := decimal.NewFromFloat(floatA).Round(place)
	deciB := decimal.NewFromFloat(floatB).Round(place)
	if deciA.GreaterThan(deciB) {
		return 1
	} else if deciA.LessThan(deciB) {
		return -1
	}
	return 0
}
