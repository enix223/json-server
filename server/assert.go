package server

func Assert(cond bool, message string) {
	if !cond {
		panic(message)
	}
}
