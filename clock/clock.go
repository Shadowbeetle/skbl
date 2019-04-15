package clock

import "time"

type Timer interface {
	Reset(time.Duration) bool
}
