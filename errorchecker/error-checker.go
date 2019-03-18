package errorchecker

type ErrorChecker struct {
	done chan bool
}

func NewErrorChecker(doneChannel chan bool) ErrorChecker {
	return ErrorChecker{
		done: doneChannel,
	}
}

func (ec *ErrorChecker) Check(e error) {
	if e != nil {
		close(ec.done)
		panic(e)
	}
}
