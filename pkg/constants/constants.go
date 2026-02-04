package Constants

// Обычные константы, путь до файлов или сообщения для логгера (для логерра info и warning сообщения)
const (
	LoggerPathLinux   = "./Log.txt"
	LoggerPathWindows = "../../Log.txt"
	LoggerPathDarvin  = "../../Log.txt"
)

// Константы для ошибок
const (
	LogFileDoesNotOpen  = "LogFile doesn't open"
	LogFileDoesNotWrite = "LogFile doesn't write"
)
