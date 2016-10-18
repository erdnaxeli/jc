package dispatcher

import (
	"log"

	"git.iiens.net/morignot2011/jc"
)

func (d *Dispatcher) join(transport string, ev *jc.JoinEvent) {
	log.Printf("Receive JoinEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	links := d.findLink(transport, ev.Channel)
	for _, link := range links {
		if isFiltered(link.filters, ev.Nick) {
			continue
		}

		for _, endpoint := range link.endpoints {
			if endpoint.transport == transport && endpoint.channel == ev.Channel {
				continue
			}

			d.transports[endpoint.transport].Join(&jc.JoinEvent{
				Nick:    ev.Nick + "_jc",
				Channel: endpoint.channel,
			})
		}
	}
}

func (d *Dispatcher) message(transport string, ev *jc.MessageEvent) {
	log.Printf("Receive MessageEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	links := d.findLink(transport, ev.Channel)
	for _, link := range links {
		if isFiltered(link.filters, ev.Nick) {
			continue
		}

		for _, endpoint := range link.endpoints {
			if endpoint.transport == transport && endpoint.channel == ev.Channel {
				continue
			}

			d.transports[endpoint.transport].Message(&jc.MessageEvent{
				Nick:    ev.Nick + "_jc",
				Channel: endpoint.channel,
				Text:    ev.Text,
			})
		}
	}
}
