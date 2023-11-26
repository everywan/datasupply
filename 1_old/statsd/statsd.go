package statsd

import (
	"time"
)

type IStatsd interface {
	Gauge(name string, value float64)
	Increment(metrics string)
	Count(name string, value int64)
	Timing(name string, value time.Duration)
	TimingUtilNow(name string, start time.Time)
	TimeInMilliseconds(name string, value float64)
	TimeInMillisecondsUtilNow(name string, start float64)
	Close() error
}

type EmptyStatsd struct{}

var _ IStatsd = new(EmptyStatsd)

func (statsd *EmptyStatsd) Gauge(name string, value float64)                     {}
func (statsd *EmptyStatsd) Increment(metrics string)                             {}
func (statsd *EmptyStatsd) Count(name string, value int64)                       {}
func (statsd *EmptyStatsd) Timing(name string, value time.Duration)              {}
func (statsd *EmptyStatsd) TimingUtilNow(name string, start time.Time)           {}
func (statsd *EmptyStatsd) TimeInMilliseconds(name string, value float64)        {}
func (statsd *EmptyStatsd) TimeInMillisecondsUtilNow(name string, start float64) {}
func (statsd *EmptyStatsd) Close() error                                         { return nil }
