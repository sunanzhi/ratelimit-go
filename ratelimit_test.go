package ratelimit

import (
	"fmt"
	"testing"
)

func BenchmarkLimiting(t *testing.B) {
	var slidewindow, _ = Init(10, 1000, 100)
	for i := 0; i < 10000; i++ {

		err := slidewindow.Limiting()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
