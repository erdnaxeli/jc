package irc

import (
	"log"

	"git.iiens.net/morignot2011/jc"
)

func (t *Transport) Join(ev *jc.JoinEvent) {
	if i := FindChannel(t.channels, ev.Channel); i == -1 {
		// cannot join a channel not configured
		return
	}

	client, ok := t.userClients[ev.Nick]
	if !ok {
		ircCfg, err := t.newIrcConfig(ev.Nick, t.cfg)
		if err != nil {
			log.Print(err)
			return
		}

		client = t.getIrcClient(ircCfg)

		t.userClients[ev.Nick] = client
		t.userChannels[ev.Nick] = []string{ev.Channel}

		if err := client.Connect(); err != nil {
			log.Print(err)
			return
		}
	} else {
		// add channel to the user list
		if i := FindChannel(t.userChannels[ev.Nick], ev.Channel); i == -1 {
			t.userChannels[ev.Nick] = append(t.userChannels[ev.Nick], ev.Channel)
			client.Join(ev.Channel)
		}
	}
}

func (t *Transport) Message(ev *jc.MessageEvent) {
	client, ok := t.userClients[ev.Nick]
	if !ok {
		// unknown client
		return
	}

	if i := FindChannel(t.userChannels[ev.Nick], ev.Channel); i == -1 {
		// this user is not on this channel
		return
	}

	client.Privmsg(ev.Channel, ev.Text)
}

func (t *Transport) Nick(ev *jc.NickEvent) {
	client, ok := t.userClients[ev.OldNick]
	if !ok {
		// unknown client
		return
	}

	t.realNicks[client.Me().Nick] = ev.NewNick
}

func (t *Transport) PrivMessage(ev *jc.PrivMessageEvent) {
	client, ok := t.userClients[ev.Nick]
	if !ok {
		// unknown client
		return
	}

	client.Privmsg(ev.Channel, ev.Text)
}

func (t *Transport) Part(ev *jc.PartEvent) {
	client, ok := t.userClients[ev.Nick]
	if !ok {
		// unknown client
		return
	}

	i := FindChannel(t.userChannels[ev.Nick], ev.Channel)
	if i := FindChannel(t.userChannels[ev.Nick], ev.Channel); i == -1 {
		// this user is not on this channel
		return
	}

	client.Part(ev.Channel)
	// remove chan from user's
	t.userChannels[ev.Nick] = append(t.userChannels[ev.Nick][:i], t.userChannels[ev.Nick][i+1:]...)

	// make the user quit if is not in anymore channel
	if len(t.userChannels) == 0 {
		t.Quit(&jc.QuitEvent{ev.Nick})
	}
}

func (t *Transport) Quit(ev *jc.QuitEvent) {
	client, ok := t.userClients[ev.Nick]
	if !ok {
		// unknown client
		return
	}

	client.Quit()
	delete(t.userClients, ev.Nick)
	delete(t.userChannels, ev.Nick)
	delete(t.realNicks, client.Me().Nick)
}
