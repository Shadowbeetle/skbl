package backlight

type timerError struct {
	s string
}

var TimerError timerError = timerError{"kbl.timer is not time.Timer"}

func (e timerError) Error() string {
	return e.s
}
