package speedrun

import (
	"fmt"
	"time"
)

// Stopwatch represents a timer on top of a game.
// It can handle start, stop, pause/unpause and reset actions.
type Stopwatch struct {
	pauseTime          time.Time
	startTime          time.Time
	stopTime           time.Time
	pauseTotalDuration time.Duration
}

// Start resets the stopwatch and start it.
func (s *Stopwatch) Start() {
	s.Reset()
	s.startTime = time.Now()
}

// Pause function pause or unpause the stopwatch depending on its state.
// If the stopwatch has been stopped previously, this function has no effect.
func (s *Stopwatch) Pause() {
	// Cannot pause a stopped or unstarted timer
	if !s.stopTime.IsZero() || s.startTime.IsZero() {
		return
	}

	if s.pauseTime.IsZero() {
		s.pauseTime = time.Now()
	} else {
		s.pauseTotalDuration += time.Since(s.pauseTime)
		s.pauseTime = time.Time{}
	}
}

// Stop halt the stopwatch and cannot be resumed again
func (s *Stopwatch) Stop() {
	s.stopTime = time.Now()
}

// Reset set the stopwatch to all its default value
func (s *Stopwatch) Reset() {
	s.pauseTotalDuration = time.Duration(0)
	s.pauseTime = time.Time{}
	s.stopTime = time.Time{}
	s.startTime = time.Time{}
}

// Elapsed return the duration since the stopwatch has been started.
// Pause time are removed from elapsed time.
func (s *Stopwatch) Elapsed() time.Duration {
	// If not started, return 0
	if s.startTime.IsZero() {
		return time.Duration(0)
	}

	refTime := time.Now()
	// If stopped or pause, current reference time is the time
	// since the stopwatch is in its state
	if !s.stopTime.IsZero() {
		refTime = s.stopTime
	} else if !s.pauseTime.IsZero() {
		refTime = s.pauseTime
	}

	// Paused duration should not be taken into account in elapsed time
	return refTime.Sub(s.startTime) - s.pauseTotalDuration
}

// StopwatchFormatHour format a duration to a stopwatch string format for hour scale.
// Format: hh:mm:ss.dth
func StopwatchFormatHour(d time.Duration) string {
	hundredth := int64(d.Milliseconds()/100) % 10
	seconds := int64(d.Seconds()) % 60
	minutes := int64(d.Minutes()) % 60
	hours := int64(d.Hours())

	return fmt.Sprintf("%02d:%02d:%02d.%01d", hours, minutes, seconds, hundredth)
}

// StopwatchFormatMinute format a duration to a minutes stopwatch string format for minutes scale.
// Format: mm:ss.hth
func StopwatchFormatMinute(d time.Duration) string {
	hundredth := int64(d.Milliseconds()/10) % 100
	seconds := int64(d.Seconds()) % 60
	minutes := int64(d.Minutes())

	return fmt.Sprintf("%02d:%02d.%02d", minutes, seconds, hundredth)
}
