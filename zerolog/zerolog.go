package zerolog

import (
	"fmt"
	"strings"
	"time"
	// "log"
)

type Writer interface {
	// Write(loginId string, ticketingId string, stepName string, message string)
	Write(appName string, goodsCode string, loginId string, ticketingId string, stepName string, message string)
}

type Event struct {
	// contains filtered or unexported fields
	l         *Logger
	StepValue string
	MsgValue  string
}

func (e *Event) Step(step string) *Event {
	e.StepValue = step
	return e
}

func (e *Event) Msg(msg string) {
	e.MsgValue = msg
	// println(msg)
	e.l.w.Write(e.l.c.AppName, e.l.c.GoodsCode, e.l.c.LoginId, e.l.c.TkId, e.StepValue, e.MsgValue)
}

func (e *Event) MsgArr(a ...interface{}) {

	var b strings.Builder
	for _, item := range a {

		if w, ok := item.(string); ok {
			fmt.Fprintf(&b, "%s | ", w)
		} else if w, ok := item.(time.Time); ok {
			fmt.Fprintf(&b, "%s | ", w.Format("15:04:05.000"))
		} else if w, ok := item.(time.Duration); ok {
			fmt.Fprintf(&b, "%.2fms | ", float64(w.Microseconds())/1000.00)
		} else if w, ok := item.(int); ok {
			fmt.Fprintf(&b, "%d | ", w)
		} else {
			fmt.Fprintf(&b, "%+v | ", item)
		}

	}

	e.MsgValue = strings.TrimSuffix(b.String(), " | ")
	e.l.w.Write(e.l.c.AppName, e.l.c.GoodsCode, e.l.c.LoginId, e.l.c.TkId, e.StepValue, e.MsgValue)
}

// func (e *Event) Str(key, val string) *Event {
// 	return nil
// }
// func (e *Event) Stack() *Event {
// 	return nil
// }
// func (e *Event) Err(err error) *Event {
// 	return nil
// }

// func (e *Event) Dur(key string, d time.Duration) *Event {
// 	return nil
// }
// func (e *Event) Time(key string, t time.Time) *Event {
// 	return nil
// }
// func (e *Event) Int(key string, i int) *Event {
// 	return nil
// }

// var Logger = zerolog.New(StderrWriter{}).With().Logger()
func New(w Writer) Logger {
	return Logger{w, nil, nil}
}

type Logger struct {
	w Writer
	c *Context
	e *Event
}

func (l Logger) With() Context {
	l.c = &Context{&l, "", "", "", ""}
	return *l.c
}

func (l *Logger) Printf(format string, v ...interface{}) {
}

func (l *Logger) Log() *Event {
	l.e = &Event{l, "", ""}
	return l.e
}

// func (l *Logger) Error() *Event {
// 	return nil
// }

type Context struct {
	// contains filtered or unexported fields
	l         *Logger
	AppName   string
	GoodsCode string
	LoginId   string
	TkId      string
}

func (c Context) Logger() Logger {
	c.l.c = &c
	return *c.l
}

func (c Context) SetAppName(appName string) Context {
	c.AppName = appName
	return c
}

func (c Context) SetGoodsCode(goodsCode string) Context {
	c.GoodsCode = goodsCode
	return c
}

func (c Context) SetLoginId(loginId string) Context {
	c.LoginId = loginId
	return c
}

func (c Context) SetTkId(tkId string) Context {
	c.TkId = tkId
	return c
}

// func LogTkEvent(tkinfo TicketingInfo, stepName string, message string) {
//
// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				log.Log().Str("f", fmt.Sprintf("%v", r)).Msg("Error is recovered.")
// 				var err error
// 				switch x := r.(type) {
// 				case string:
// 					err = errors.New(x)
// 				case error:
// 					err = x
// 				default:
// 					err = errors.New("unknown panic")
// 				}
// 				log.Error().Stack().Err(err).Msg("")
// 			}
// 		}()
//
// 		// set up a connection to the server.
// 		conn, err := grpc.Dial(cfg.LoggerAddress, grpc.WithInsecure(), grpc.WithBlock())
// 		if err != nil {
// 			log.Error().Err(err).Msg("")
// 		}
// 		defer conn.Close()
//
// 		c := pb.NewGreeterClient(conn)
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 		defer cancel()
// 		{
// 			_, err := c.LogTicketingEvent(ctx, &pb.TicketingEventRequest{LoginId: tkinfo.SiteLoginId, TkId: tkinfo.TkId, StepName: stepName, Message: message})
// 			if err != nil {
// 				log.Error().Err(err).Msg("")
// 			}
// 		}
//
// 	}()
// }

// func SetChanSize(size int) {
// 	logchan = make(chan string, size)
// }
//
// func SetFlags(flag int) {
// 	log.SetFlags(flag)
// }
//
// func Print(v ...interface{}) {
// 	// log.Printf(v...)
// 	logchan <- fmt.Sprint(v...)
// }
//
// func Printf(format string, v ...interface{}) {
// 	// log.Printf(format, v...)
// 	logchan <- fmt.Sprintf(format, v...)
// }
//
// func Println(a ...interface{}) (n int, err error) {
// 	logchan <- fmt.Sprintln(a...)
// 	return
// }
