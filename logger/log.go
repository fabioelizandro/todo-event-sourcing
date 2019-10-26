package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogMsg interface {
	FieldStr(name string, value string) LogMsg
	FieldErr(err error) LogMsg
	FieldInterface(name string, value interface{}) LogMsg
	Write()
}

type Log interface {
	InfoMsg(msg string) LogMsg
	ErrorMsg(msg string) LogMsg
}

type noopLog struct {
}

func (n *noopLog) InfoMsg(msg string) LogMsg {
	return &noopLogMsg{}
}

func (n *noopLog) ErrorMsg(msg string) LogMsg {
	return &noopLogMsg{}
}

type noopLogMsg struct {
}

func (n *noopLogMsg) FieldInterface(string, interface{}) LogMsg {
	return n
}

func (n *noopLogMsg) FieldErr(err error) LogMsg {
	return n
}

func (n *noopLogMsg) FieldStr(name string, value string) LogMsg {
	return n
}

func (n *noopLogMsg) Write() {

}

type zLog struct {
}

func (z *zLog) InfoMsg(msg string) LogMsg {
	return &zLogMsg{msg: msg, zEvent: log.Info()}
}

func (z *zLog) ErrorMsg(msg string) LogMsg {
	return &zLogMsg{msg: msg, zEvent: log.Error()}
}

type zLogMsg struct {
	msg    string
	zEvent *zerolog.Event
}

func (z *zLogMsg) FieldInterface(name string, value interface{}) LogMsg {
	return &zLogMsg{msg: z.msg, zEvent: z.zEvent.Interface(name, value)}
}

func (z *zLogMsg) FieldErr(err error) LogMsg {
	return &zLogMsg{msg: z.msg, zEvent: z.zEvent.Err(err)}
}

func (z *zLogMsg) FieldStr(name string, value string) LogMsg {
	return &zLogMsg{msg: z.msg, zEvent: z.zEvent.Str(name, value)}
}

func (z *zLogMsg) Write() {
	z.zEvent.Msg(z.msg)
}

func NewZLog() *zLog {
	return &zLog{}
}

func NewNoopLog() *noopLog {
	return &noopLog{}
}
