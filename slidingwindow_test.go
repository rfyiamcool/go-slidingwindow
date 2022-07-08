package slidingwindow

import (
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	sw := New(100, 1*time.Second)
	due := time.Now().Add(5 * time.Second)
	go func() {
		for {
			if time.Now().After(due) {
				return
			}

			sw.Set(1)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	for {
		if time.Now().After(due) {
			return
		}

		t.Log("100ms sum", sw.GetStatsWithDuration(100*time.Millisecond).Sum)
		t.Log("200ms sum", sw.GetStatsWithDuration(200*time.Millisecond).Sum)
		time.Sleep(100 * time.Millisecond)
	}
}
