package rabbit

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Error(string)
	Errorf(string, ...any)
}
