package irc

import (
	//"fmt"
	"log"
	"strings"
	//"time"

	irc "github.com/fluffle/goirc/client"

	"git.iiens.net/morignot2011/jc"
)

func (t *Transport) connected(client *irc.Conn, line *irc.Line) {
	t.Logger.Printf("%s is connected", client.Me().Nick)
	if t.client != client {
		// user's client
		for _, channel := range t.userChannels[client.Me().Nick] {
			log.Printf("Join %s", channel)
			client.Join(channel)
		}
	} else {
		// bot's client
		//st := t.client.StateTracker()
		//me := st.Me()
		for _, channel := range t.channels {
			t.client.Join(channel)
			/*
				time.Sleep(time.Second * 2)

				if _, ok := me.Channels[channel]; !ok {
					t.connectionError <- fmt.Errorf("Cannot join channel %s", channel)
					break
				}
			*/
		}

		close(t.connectionError)
	}
}

func (t *Transport) disconnected(client *irc.Conn, line *irc.Line) {
	if t.client != client {
		// all should already have been cleaned
		log.Printf("%s got disconnected", client.Me().Nick)
		return
	}

	// bot's client
	t.End <- true
}

func (t *Transport) join(client *irc.Conn, line *irc.Line) {
	if t.client != client || t.isUserDistant(line.Nick) {
		return
	}

	me := t.client.StateTracker().Me()
	if line.Nick == me.Nick {
		return
	}

	// bot's client
	t.Events <- &jc.JoinEvent{
		Nick:    line.Nick,
		Channel: line.Args[0],
	}
}

func (t *Transport) nick(client *irc.Conn, line *irc.Line) {
	if t.client != client || t.isUserDistant(line.Nick) {
		return
	}

	// bot's client
	t.Events <- &jc.NickEvent{
		OldNick: line.Nick,
		NewNick: line.Args[0],
	}
}

func (t *Transport) part(client *irc.Conn, line *irc.Line) {
	if t.client != client || t.isUserDistant(line.Nick) {
		return
	}

	// bot's client
	for _, channel := range line.Args {
		t.Events <- &jc.PartEvent{
			Nick:    line.Nick,
			Channel: channel,
		}
	}
}

func (t *Transport) privmsg(client *irc.Conn, line *irc.Line) {
	text := line.Args[len(line.Args)-1]

	if t.client != client {
		for _, target := range line.Args[:len(line.Args)-1] {
			if target[0] == '#' || target[0] == '$' {
				continue
			}

			// this is a query
			t.Events <- &jc.PrivMessageEvent{
				Nick:    line.Nick,
				Channel: target,
				Text:    text,
			}
		}
	} else if !t.isUserDistant(line.Nick) {
		// bot's client
		for _, target := range line.Args[:len(line.Args)-1] {
			if target[0] != '#' || strings.Contains(target, ".") {
				// Something else than # is either a server mask or a query
				// # with a . in the target is a host mask
				continue
			}

			t.Events <- &jc.MessageEvent{
				Nick:    line.Nick,
				Channel: target,
				Text:    text,
			}
		}
	}
}

func (t *Transport) quit(client *irc.Conn, line *irc.Line) {
	if t.client != client || t.isUserDistant(line.Nick) {
		return
	}

	// bot's client
	t.Events <- &jc.QuitEvent{
		Nick: line.Nick,
	}
}
