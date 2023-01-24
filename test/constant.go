package test

import "time"

var (
	TimeBeforeNow = time.Now().Add(-1 * time.Hour).UTC().Round(time.Microsecond)
	TimeAfterNow  = time.Now().Add(1 * time.Hour).UTC().Round(time.Microsecond)
)
