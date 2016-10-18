package dispatcher

import (
	"log"
	"reflect"

	"git.iiens.net/morignot2011/jc"
)

type Dispatcher struct {
	links          []Link
	transports     map[string]jc.Transport
	transportNames []string
}

type Link struct {
	endpoints []Endpoint
	filters   []string
}

type Endpoint struct {
	channel   string
	transport string
}

func NewDispatcher(cfg []map[string]interface{}, transports map[string]jc.Transport) jc.Dispatcher {
	d := &Dispatcher{
		transports: transports,
	}

	for _, linkCfg := range cfg {
		link := Link{}

		filters, ok := linkCfg["filters"]
		if ok {
			for _, filter := range filters.([]interface{}) {
				link.filters = append(link.filters, filter.(string))
			}
		}

		delete(linkCfg, "filters")

		for name, channels := range linkCfg {
			for _, channel := range channels.([]interface{}) {
				endpoint := Endpoint{
					channel:   channel.(string),
					transport: name,
				}

				link.endpoints = append(link.endpoints, endpoint)
			}
		}

		d.links = append(d.links, link)
	}

	return d
}

func (d *Dispatcher) Run() {
	var cases []reflect.SelectCase
	var transportNames []string

	for name, t := range d.transports {
		cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.GetEvents())})
		transportNames = append(transportNames, name)
	}

	for {
		i, value, ok := reflect.Select(cases)
		name := transportNames[i]

		if !ok {
			log.Fatalf("Transport %s closed its event channel!", name)
		}

		switch ev := value.Interface().(type) {
		case *jc.JoinEvent:
			d.join(name, ev)
		case *jc.MessageEvent:
			d.message(name, ev)
		}
	}
}

func (d *Dispatcher) findLink(transport string, channels ...string) []Link {
	byChan := false
	if len(channels) > 0 {
		byChan = true
	}

	var links []Link

	for _, link := range d.links {
		match := false

		for _, endpoint := range link.endpoints {
			if endpoint.transport != transport {
				continue
			}

			if byChan && endpoint.channel != channels[0] {
				continue
			}

			match = true
			break
		}

		if match {
			links = append(links, link)
		}
	}

	return links
}

func isFiltered(filters []string, nick string) bool {
	for _, v := range filters {
		if v == nick {
			return true
		}
	}

	return false
}
