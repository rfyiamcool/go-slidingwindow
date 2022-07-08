package slidingwindow

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindow struct {
	ring         []*WindowEntry
	idx          int
	slots        int
	slotInterval time.Duration
	duration     time.Duration

	sync.RWMutex
}

type WindowEntry struct {
	Datetime time.Time
	Cnt      int64
}

func (we *WindowEntry) Clone() *WindowEntry {
	return &WindowEntry{
		Datetime: we.Datetime,
		Cnt:      we.Cnt,
	}
}

func NewSlidingWindow(slots int, duration time.Duration) *SlidingWindow {
	var (
		ring []*WindowEntry
	)

	for i := 0; i < slots; i++ {
		ring = append(ring, &WindowEntry{})
	}

	sw := &SlidingWindow{
		slots:        slots,
		duration:     duration,
		ring:         ring,
		slotInterval: duration / time.Duration(slots),
	}

	return sw
}

func (s *SlidingWindow) Dump() []*WindowEntry {
	s.RLock()
	defer s.RUnlock()

	res := []*WindowEntry{}
	for _, entry := range s.ring {
		s.moveIndex()
		if entry.Datetime.IsZero() {
			continue
		}
		res = append(res, entry.Clone())
	}
	return res
}

func (s *SlidingWindow) Set(cnt int64) {
	s.Lock()
	defer s.Unlock()

	s.set(time.Now(), cnt)
}

func (s *SlidingWindow) SetWithTime(tm time.Time, cnt int64) {
	s.Lock()
	defer s.Unlock()

	s.set(tm, cnt)
}

func (s *SlidingWindow) set(now time.Time, cnt int64) {
	due := s.ring[s.idx].Datetime.Add(s.slotInterval)
	if now.After(due) {
		s.moveIndex()
		s.ring[s.idx].Datetime = now
		s.ring[s.idx].Cnt = cnt
		return
	}

	s.ring[s.idx].Cnt += cnt
}

func (s *SlidingWindow) GetSum() int64 {
	return s.GetStatsWithDuration(s.duration).Sum
}

func (s *SlidingWindow) GetAvg() int64 {
	return s.GetStatsWithDuration(s.duration).Avg
}

func (s *SlidingWindow) GetMin() int64 {
	return s.GetStatsWithDuration(s.duration).Min
}

func (s *SlidingWindow) GetMax() int64 {
	return s.GetStatsWithDuration(s.duration).Max
}

func (s *SlidingWindow) GetSumWithDuration(dur time.Duration) int64 {
	return s.GetStatsWithDuration(dur).Sum
}

func (s *SlidingWindow) GetAvgWithDuration(dur time.Duration) int64 {
	return s.GetStatsWithDuration(dur).Avg
}

func (s *SlidingWindow) GetMinWithDuration(dur time.Duration) int64 {
	return s.GetStatsWithDuration(dur).Min
}

func (s *SlidingWindow) GetMaxWithDuration(dur time.Duration) int64 {
	return s.GetStatsWithDuration(dur).Max
}

type StatsEntry struct {
	Sum int64
	Avg int64
	Min int64
	Max int64
}

func (se StatsEntry) String() string {
	return fmt.Sprintf("sum: %v, avg: %v, min: %v, max: %v",
		se.Sum, se.Avg, se.Min, se.Max,
	)
}

func (s *SlidingWindow) GetStats() StatsEntry {
	s.Lock()
	defer s.Unlock()

	return s.getStats(s.duration)
}

func (s *SlidingWindow) GetStatsWithDuration(dur time.Duration) StatsEntry {
	s.Lock()
	defer s.Unlock()

	return s.getStats(dur)
}

func (s *SlidingWindow) getStats(duration time.Duration) StatsEntry {
	var (
		sumv int64
		minv int64
		maxv int64
	)

	point := time.Now().Add(-duration) // last interval
	for _, entry := range s.ring {
		if entry.Datetime.Before(point) {
			continue
		}

		if entry.Cnt < minv {
			minv = entry.Cnt
		}
		if entry.Cnt < maxv {
			maxv = entry.Cnt
		}
		sumv += entry.Cnt
	}

	return StatsEntry{
		Sum: sumv,
		Min: minv,
		Max: maxv,
		Avg: sumv / int64(len(s.ring)),
	}
}

func (s *SlidingWindow) Reset() {
	s.Lock()
	defer s.Unlock()

	for _, en := range s.ring {
		en.Cnt = 0
	}
}

func (s *SlidingWindow) ResetHead() {
	s.Lock()
	defer s.Unlock()

	// lookup head offset
	idx := len(s.ring) - s.idx

	// abs
	if idx < 0 {
		idx = 0 - idx
	}
	s.ring[idx].Cnt = 0
}

// moveIndex
func (s *SlidingWindow) moveIndex() {
	s.idx++
	if s.idx >= s.slots {
		s.idx = 0
	}
}
