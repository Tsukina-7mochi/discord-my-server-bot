package botlog

import (
	"fmt"
	"log"
	. "mochi-bot/internal/pkg/discordchannelwriter"

	"github.com/bwmarrin/discordgo"
)

const (
	Debug = "DEBUG"
	Info  = "INFO"
	Warn  = "WARN"
	Error = "ERROR"
)

type Logger struct {
	defaultLogger *log.Logger
	session       *discordgo.Session
	channelWriter *DiscordChannelWriter
}

func NewLogger(defaultLogger *log.Logger, session *discordgo.Session, channelID string) *Logger {
	return &Logger{
		defaultLogger: defaultLogger,
		session:       session,
		channelWriter: NewDiscordChannelWriter(session, channelID),
	}
}

func (l *Logger) writef(level string, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.defaultLogger.Println(level + ": " + message)

	_, err := l.channelWriter.Write([]byte("`" + level + "`: " + message))
	if err != nil {
		l.defaultLogger.Printf("%s: Failed to write to Discord channel: %v", Warn, err)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.writef(Debug, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.writef(Info, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.writef(Warn, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.writef(Error, format, args...)
}

func (l *Logger) CommandDebugf(commandName string, format string, args ...interface{}) {
	l.Debugf("[command `/%s`] "+format, append([]interface{}{commandName}, args...)...)
}

func (l *Logger) CommandInfof(commandName string, format string, args ...interface{}) {
	l.Infof("[command `/%s`] "+format, append([]interface{}{commandName}, args...)...)
}

func (l *Logger) CommandWarnf(commandName string, format string, args ...interface{}) {
	l.Warnf("[command `/%s`] "+format, append([]interface{}{commandName}, args...)...)
}

func (l *Logger) CommandErrorf(commandName string, format string, args ...interface{}) {
	l.Errorf("[command `/%s`] "+format, append([]interface{}{commandName}, args...)...)
}
