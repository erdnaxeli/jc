package jc

import (
	"log"
)

type Transport interface {
	Run() (err error)

	Join(*JoinEvent)
	Message(*MessageEvent)
	Nick(*NickEvent)
	PrivMessage(*PrivMessageEvent)
	Part(*PartEvent)
	Quit(*QuitEvent)

	GetEvents() chan interface{}
	GetEnd() chan bool
}

type BaseTransport struct {
	Events chan interface{}
	End    chan bool

	Logger *log.Logger
}

func (t *BaseTransport) GetEvents() chan interface{} {
	return t.Events
}

func (t *BaseTransport) GetEnd() chan bool {
	return t.End
}
