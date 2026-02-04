package logger

type MyLogger interface {
	Info(string, string)
	Warning(string, string)
	Error(string, string)
}
