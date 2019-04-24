package backlight

type timerError struct {
	s string
}

var TimerError timerError = timerError{"timer is not time.Timer"}
var TimerConfigError timerError = timerError{"do not set timerCh without setting timer first"}

func (e timerError) Error() string {
	return e.s
}
