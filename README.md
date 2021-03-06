# slidingwindow

golang simple slidingwindow lib

### api

- New(slots int, duration time.Duration)
- Set(cnt int64)
- SetWithTIme(tm time.Time, cnt int64)
- GetSum()
- GetAvg()
- GetMin()
- GetMax()
- GetSumWithDuration(dur time.Duration)
- GetAvgWithDuration(dur time.Duration)
- GetMinWithDuration(dur time.Duration)
- GetMaxWithDuration(dur time.Duration)
- GetStats()
- GetStatsWithDuration(dur time.Duration)

### example

```go
package main

import (
	"fmt"
	"time"

	"github.com/rfyiamcool/go-slidingwindow"
)

func main() {
	sw := slidingwindow.New(100, 1*time.Second)
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

		fmt.Println("-1s    stats", sw.GetStats().String())
		fmt.Println("-100ms sum", sw.GetStatsWithDuration(100*time.Millisecond).Sum)
		fmt.Println("-200ms max", sw.GetStatsWithDuration(200*time.Millisecond).Max)
		time.Sleep(100 * time.Millisecond)
	}
}
```
