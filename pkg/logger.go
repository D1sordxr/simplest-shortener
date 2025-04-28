package pkg

type Log interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}
