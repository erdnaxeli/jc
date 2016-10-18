package jc

type Dispatcher interface {
	Run()
}

type Links []Link

type Link struct {
	endpoints []Endpoint
	filters   []string
}

type Endpoint struct {
	Channel   string
	Transport string
}
