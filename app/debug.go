package app

func (a *APITest) iStartDebugging() {
	a.debug = true
}

func (a *APITest) iStopDebugging() {
	a.debug = false
}
