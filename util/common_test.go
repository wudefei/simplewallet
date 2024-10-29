package util_test

import (
	"simplewallet/util"
	"testing"

	"go.uber.org/goleak"
)

func TestCompareFloat(t *testing.T) {
	defer goleak.VerifyNone(t) // check for goroutine leaks

	type args struct {
		floatA float64
		floatB float64
		place  int32
	}
	type caseInfo struct {
		Name string
		Args args
		Want int
	}
	tests := []caseInfo{
		{
			Name: "case1: floatA == floatB",
			Args: args{floatA: 1.111111, floatB: 1.111111, place: 6}, // 1.111111 == 1.111111
			Want: 0,
		},
		{
			Name: "case2: floatA > floatB",
			Args: args{floatA: 1.111111, floatB: 1.111110, place: 6}, // 1.111111 > 1.111110
			Want: 1,
		},
		{
			Name: "case3: floatA < floatB",
			Args: args{floatA: 1.111110, floatB: 1.111111, place: 6}, // 1.111110 < 1.111111
			Want: -1,
		},
		{
			Name: "case4: floatA > 0, floatB < 0",
			Args: args{floatA: 1.111111, floatB: -1.111111, place: 6}, // 1.111111 > -1.111111
			Want: 1,
		},
		{
			Name: "case5: floatA = 0, floatB < 0",
			Args: args{floatA: 0, floatB: -1.111111, place: 6}, // 0 > -1.111111
			Want: 1,
		},
		{
			Name: "case6: floatA = 0, floatB > 0",
			Args: args{floatA: 0, floatB: 1.111111, place: 6}, // 0 < 1.111111
			Want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if got := util.CompareFloat(tt.Args.floatA, tt.Args.floatB, tt.Args.place); got != tt.Want {
				t.Errorf("CompareFloat() = %v, want %v", got, tt.Want)
			}
		})
	}
}
