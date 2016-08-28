package metrics

type errString struct {
	msg string
}

func (e errString) Error() string {
	return e.msg
}

// InitError indicates that the env hasn't been setup right
type InitError struct {
	errString
}
