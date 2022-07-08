# slidingwindow

golang simple slidingwindow lib

### api

- NewSlidingWindow(slots int, duration time.Duration)
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
sw := NewSlidingWindow(30, 1*time.Minute)

sw.Set(1)
time.Sleep(1000 * time.Millisecond)

sw.Set(1)
time.Sleep(2000 * time.Millisecond)

fmt.Println("-1min sum", sw.GetStats().Sum)
fmt.Println("-10s  sum", sw.GetStatsWithDuration(10*time.Second).Sum)
```
