package dispatcher

import (
	"log"

	"git.iiens.net/morignot2011/jc"
)

func (d *Dispatcher) join(transport string, ev *jc.JoinEvent) {
	log.Printf("Receive JoinEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	endpoints := d.findTransports(transport, ev.Nick, ev.Channel)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].Join(&jc.JoinEvent{
			Nick:    ev.Nick + "_jc",
			Channel: endpoint.channel,
		})
	}
}

func (d *Dispatcher) message(transport string, ev *jc.MessageEvent) {
	log.Printf("Receive MessageEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	endpoints := d.findTransports(transport, ev.Nick, ev.Channel)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].Message(&jc.MessageEvent{
			Nick:    ev.Nick + "_jc",
			Channel: endpoint.channel,
			Text:    ev.Text,
		})
	}
}

func (d *Dispatcher) nick(transport string, ev *jc.NickEvent) {
	log.Printf("Receive NickEvent from %s : %s to %s", transport, ev.OldNick, ev.NewNick)

	endpoints := d.findTransports(transport, ev.OldNick)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].Nick(&jc.NickEvent{
			OldNick: ev.OldNick,
			NewNick: ev.NewNick,
		})
	}
}

func (d *Dispatcher) privMessage(transport string, ev *jc.PrivMessageEvent) {
	log.Printf("Receive PrivMessageEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	endpoints := d.findTransports(transport, ev.Nick, ev.Channel)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].PrivMessage(&jc.PrivMessageEvent{
			Nick:    ev.Nick + "_jc",
			Channel: endpoint.channel,
			Text:    ev.Text,
		})
	}
}

func (d *Dispatcher) part(transport string, ev *jc.PartEvent) {
	log.Printf("Receive PartEvent from %s: %s on %s", transport, ev.Nick, ev.Channel)

	endpoints := d.findTransports(transport, ev.Nick, ev.Channel)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].Part(&jc.PartEvent{
			Nick:    ev.Nick + "_jc",
			Channel: endpoint.channel,
		})
	}
}

func (d *Dispatcher) quit(transport string, ev *jc.QuitEvent) {
	log.Printf("Receive QuitEvent from %s: %s", transport, ev.Nick)

	endpoints := d.findTransports(transport, ev.Nick)
	for _, endpoint := range endpoints {
		d.transports[endpoint.transport].Quit(&jc.QuitEvent{
			Nick: ev.Nick + "_jc",
		})
	}
}
