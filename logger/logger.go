package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type (
	// Log é o modelo default para logs de TagManager
	Log interface {
		printer(s string, v ...interface{})

		Debug(s string, v ...interface{})
		Info(s string, v ...interface{})
		Warn(s string, v ...interface{})
		Error(s string, v ...interface{})
		Critical(s string, v ...interface{})
	}

	// config é a estrutura de configuração do Logger
	config struct {
		Output io.Writer
	}

	// Logger é a implementação padrão de Log
	Logger struct {
		LogLevel int
		Config   config
	}
)

// NewLogger cria um Logger padrão para o pacote tag_manager
func NewLogger() *Logger {
	return &Logger{
		Config: config{Output: os.Stdout},
	}
}

// printer é o método que escreve no Output de config
func (l *Logger) printer(s string, v ...interface{}) {
	// Reset de loglevel
	defer func() {
		l.LogLevel = 1
	}()

	if len(v) == 0 {
		s = s + "\n"
	}

	label := "DEBUG"
	color := FgGreen

	switch l.LogLevel {
	case 2:
		label = "INFO"
		color = FgHiCyan
	case 3:
		label = "WARN"
		color = FgHiYellow
	case 4:
		label = "ERROR"
		color = FgRed
	case 5:
		label = "CRITICAL"
		color = BgRed
	default:
	}

	s = fmt.Sprintf("%s %s", Colorize(fmt.Sprintf("[%s]", label), color), s)

	if os.Getenv("ENVIRONMENT") != "testing" && osLogLevel() <= l.LogLevel {
		_, _ = fmt.Fprintf(l.Config.Output, s, v...)
	}
}

// Debug faz o print com o level debug (1)
func (l *Logger) Debug(s string, v ...interface{}) {
	l.LogLevel = 1
	l.printer(s, v...)
}

// Info faz o print com o level info (2)
func (l *Logger) Info(s string, v ...interface{}) {
	l.LogLevel = 2
	l.printer(s, v...)
}

// Warn faz o print com o level warn (3)
func (l *Logger) Warn(s string, v ...interface{}) {
	l.LogLevel = 3
	l.printer(s, v...)
}

// Error faz o print com o level error (4)
func (l *Logger) Error(s string, v ...interface{}) {
	l.LogLevel = 4
	l.printer(s, v...)
}

// Critical faz o print com o level critical (5) e destroy o processo
func (l *Logger) Critical(s string, v ...interface{}) {
	l.LogLevel = 5
	l.printer(s, v...)
	os.Exit(1)
}

// osLogLevel captura o level de logs definido nas environments da
// aplicação, caso não esteja definido, o padrão é 4 (ERROR)
func osLogLevel() int {
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		intLvl, _ := strconv.Atoi(lvl)
		return intLvl
	}

	return 4
}
