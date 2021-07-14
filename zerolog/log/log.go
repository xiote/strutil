package log

import (
	"fmt"
	"github.com/xiote/go-utils/zerolog"
	"time"
)

type StderrWriter struct {
	//
}

func (w StderrWriter) Write(appName string, goodsCode string, loginId string, ticketingId string, stepName string, message string) {

	// fmt.Fprintln(os.Stderr, "hello world")

	go println(fmt.Sprintf("%s | %s | %s | %s | %s",
		appName,
		goodsCode,
		loginId,
		ticketingId,
		// realClock.Now().Format("0102_15:04:05.000"),
		time.Now().Format("0102_15:04:05.000"),
		stepName,
		message,
	))
}

// Logger is the global logger.
// var Logger = zerolog.New(os.Stderr).With().Logger()
var Logger = zerolog.New(StderrWriter{}).With().Logger()

// func Log() *zerolog.Event {
// 	return nil
// }
//
// func Error() *zerolog.Event {
// 	return nil
// }

func With() zerolog.Context {
	return Logger.With()
}

// func Printf(format string, v ...interface{}) {
// 	// Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
// }
