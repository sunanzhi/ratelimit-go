package ratelimit

import (
	"testing"
	"time"
)

func BenchmarkLimiting(b *testing.B) {
	var slidewindow, _ = Init(10, 1000, 10000)
	for i := 0; i < 10000; i++ {

		go func() {

			err := slidewindow.Limiting()
			if err != nil {
				// fmt.Println(err.Error())
			}
		}()
	}

	time.Sleep(5 * time.Second)
}

func TestLimiting(t *testing.T) {
	var slidewindow, _ = Init(10, 1000, 10000)
	for i := 0; i < 10000; i++ {

		go func() {

			err := slidewindow.Limiting()
			if err != nil {
				// fmt.Println(err.Error())
			}
		}()
	}

	time.Sleep(5 * time.Second)
}
