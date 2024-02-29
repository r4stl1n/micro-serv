package scaff

import "time"

type ScaffoldTask struct {
	context  any
	Interval time.Duration
	Handler  func(context any, scaffoldContext *ScaffoldContext)
	LastRun  time.Time
	Synced   bool
}
