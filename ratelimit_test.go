package ratelimit

import (
	"testing"
	"time"
)

func BenchmarkLimiting(b *testing.B) {
	var (
		limitTime   = 10
		bucketCount = 5
		limitCount  = 200
	)
	var slidewindow, _ = Init(limitTime, bucketCount, limitCount)
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
